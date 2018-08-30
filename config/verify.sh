#/usr/bin/env bash

set -o errexit

# This script is meant to be run via bazel with "bazel test //config:verify"

./config/checkconfig \
    -config-path config/config.yaml \
    -job-config-path config/jobs \
    -plugin-config config/plugins.yaml
