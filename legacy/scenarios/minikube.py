#!/usr/bin/env python

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Need to figure out why this only fails on travis
# pylint: disable=bad-continuation
"""
Executes a command, and cleans up minikube state for profile $HOSTNAME
This should be run with the minikube-in-go job type
"""

import argparse
import os
import subprocess
import sys
import socket
import time

hostname = socket.gethostname()
minikube_start_cmd = [
    "minikube",
    "start",
    "--vm-driver=kvm",
    "--kubernetes-version=%s" % os.environ["KUBERNETES_VERSION"],
    "--bootstrapper=kubeadm",
    "--memory=%s" % os.environ["MINIKUBE_MEMORY"],
    "--cpus=%s" % os.environ["MINIKUBE_CPUS"],
    "--profile=%s" % hostname,
    "--feature-gates=PersistentLocalVolumes=true",
]

minikube_ingress_cmd = [
    "minikube",
    "addons",
    "enable",
    "ingress",
    "--profile=%s" % hostname,
]

minikube_dockerenv_cmd = [
    "minikube",
    "docker-env",
    "--profile=%s" % hostname,
    "--shell=sh",
]

minikube_wait_cmd = [
    "kubectl",
    "get",
    "nodes",
]

# XXX: We need the --bootstrapper argument here so that minikube knows how to
# get the logs.
# See https://github.com/kubernetes/minikube/issues/2056#issuecomment-336257971
minikube_logs_cmd = [
    "minikube",
    "--bootstrapper=kubeadm",
    "--profile", hostname,
    "logs",
]

minikube_delete_cmd = [
    "minikube",
    "delete",
    "--profile=%s" % hostname,
]

docker_ps_cmd = ["docker", "ps"]

WORKSPACE_ENV = 'WORKSPACE'
ARTIFACTS_DIRECTORY_NAME = '_artifacts'


def ensure_artifacts_directory():
    """
    Create and return the path to an artifacts directory if it doesn't already
    exist.
    """
    print >> sys.stderr, "Creating artifacts directory..."
    artifacts_path = os.path.join(
        os.getenv(WORKSPACE_ENV, os.getcwd()),
        ARTIFACTS_DIRECTORY_NAME,
    )
    try:
        os.makedirs(artifacts_path)
    except os.error as e:
        print >> sys.stderr, e
    return artifacts_path


def log_and_closefds_for_subprocess(f, cmd, args, kwargs):
    """
    Logs the command that is being run and wraps a subprocess.c* function,
    first closing stdin and FDs other than stdout and stderr, to prevent the
    subprocess doing attempting to do anything interactive. (Minikube prompts
    the user to submit a bug report if it fails.)
    """
    print >> sys.stderr, "Run: '{}'".format(" ".join(cmd))
    with open(os.devnull, "r") as devnull:
        kwargs["stdin"] = devnull
        kwargs["close_fds"] = True
        return f(cmd, *args, **kwargs)


def check_call(cmd, *args, **kwargs):
    return log_and_closefds_for_subprocess(
        subprocess.check_call,
        cmd,
        args,
        kwargs,
    )


def check_output(cmd, *args, **kwargs):
    return log_and_closefds_for_subprocess(
        subprocess.check_output,
        cmd,
        args,
        kwargs,
    )


def call(cmd, *args, **kwargs):
    return log_and_closefds_for_subprocess(
        subprocess.call,
        cmd,
        args,
        kwargs,
    )


def check(*cmd):
    """Log and run the command, raising on errors."""
    artifacts_path = ensure_artifacts_directory()
    try:
        # Run minikube start
        check_call(minikube_start_cmd)
        check_call(minikube_ingress_cmd)
        print >> sys.stderr, 'Waiting for kubernetes to become ready...'
        # Allow 2 minutes for minikube to become ready
        for i in xrange(1, 24):
            if call(minikube_wait_cmd) == 0:
                break
            time.sleep(5)
        check_call(minikube_wait_cmd)
        output = check_output(minikube_dockerenv_cmd)
        exports = output.split("\n")
        parse_exports(exports)
        check_call(docker_ps_cmd)

        print >> sys.stderr, "Execute test command"
        check_call(cmd)
    finally:
        print >> sys.stderr, 'Saving minikube logs...'
        with open(
                os.path.join(artifacts_path, "minikube-logs.txt"),
                "wb",
        ) as f:
            call(
                minikube_logs_cmd,
                stdout=f,
            )
        print >> sys.stderr, 'Deleting minikub VM...'
        call(minikube_delete_cmd)
        print >> sys.stderr, 'Deleting minikube machine files...'
        call([
            "rm", "-Rf",
            "/var/lib/libvirt/caches/minikube/.minikube/machines/%s" % hostname,
        ])


def parse_exports(exports):
    for export in exports:
        if not export.startswith("export "):
            continue
        command = export[7:].split("=")
        key = command[0]
        val = command[1]
        if val.startswith("\"") and val.endswith("\""):
            val = val[1:-1]
        os.environ[key] = val
        print >> sys.stderr, 'Setting', key, "=", val


def main(envs, cmd):
    """Run script and verify it exits 0."""
    for env in envs:
        key, val = env.split('=', 1)
        print >> sys.stderr, '%s=%s' % (key, val)
        os.environ[key] = val
    if not cmd:
        raise ValueError(cmd)
    check(*cmd)


if __name__ == '__main__':
    PARSER = argparse.ArgumentParser()
    PARSER.add_argument('--env', default=[], action='append')
    PARSER.add_argument('cmd', nargs=1)
    PARSER.add_argument('args', nargs='*')
    ARGS = PARSER.parse_args()
    main(ARGS.env, ARGS.cmd + ARGS.args)
