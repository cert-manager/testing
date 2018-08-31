# Copyright 2018 The Jetstack contributors.
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

"""Rule for building an image from a Dockerfile.
This builds the image, exports a .tar.gz and executes
container_load to make the image available to other
build rules.
"""

load("@io_bazel_rules_docker//container:container.bzl", "container_load")
load(
    "@io_bazel_rules_docker//container:pull.bzl",
    _python = "python",
)

def create(repository_ctx, tag_suffix = "", extra_args = []):
    """Core implementation of container_dockerfile."""
    build_cmd = []
    for k, v in repository_ctx.attr.build_args.items():
        build_cmd += ["--build-arg", "%s=%s" % (k, v)]

    result = repository_ctx.execute([
        repository_ctx.path(repository_ctx.attr._builder),
        "tmpfile%s" % tag_suffix, # output file name
        repository_ctx.path(repository_ctx.attr.dockerfile),
    ] + build_cmd + extra_args)

    if result.return_code:
        fail("Importing from dockerfile failed (path %s) (status %s): %s" % (repository_ctx.path(repository_ctx.attr._builder), result.return_code, result.stderr))

        """Core implementation of container_load."""

    # Add an empty top-level BUILD file.
    repository_ctx.file("BUILD", "")

    repository_ctx.file("image%s/BUILD" % tag_suffix, """
package(default_visibility = ["//visibility:public"])
load("@io_bazel_rules_docker//container:import.bzl", "container_import")
container_import(
  name = "image%s",
  config = "config.json",
  layers = glob(["*.tar"]),
)
""" % tag_suffix, executable = False)

    result = repository_ctx.execute([
        _python(repository_ctx),
        repository_ctx.path(repository_ctx.attr._importer),
        "--directory",
        repository_ctx.path("image%s" % tag_suffix),
        "--tarball",
        "tmpfile%s" % tag_suffix,
    ])

    if result.return_code:
        fail("Importing from tarball failed (status %s): %s" % (result.return_code, result.stderr))


def _impl(repository_ctx):
    variations = repository_ctx.attr.build_matrix.items()
    if len(variations) > 1:
        fail("Specifying more than one build arg variation is currently unsupported")

    varying_args = dict()
    if len(variations) == 0:
        varying_args[""] = []

    for argname, vals in variations:
        for v in vals:
            varying_args[".%s" % v] = ["--build-arg", "%s=%s" % (argname, v)]

    for suffix, args in varying_args.items():
        create(repository_ctx, suffix, args)

container_dockerfile = repository_rule(
    attrs = {
        "dockerfile": attr.label(
            allow_single_file = True,
            mandatory = True,
        ),
        "build_args": attr.string_dict(
            doc = "A dictionary of build arguments to be passed to docker build",
            allow_empty = True,
        ),
        "build_matrix": attr.string_list_dict(
            doc = "Use this parameter to build multiple variations of the same image with the given build arg set to each of the provided values. Currently only a build arg variation is supported.",
            allow_empty = True,
        ),
        "_builder": attr.label(
            executable = True,
            default = Label("//def:build-image.sh"),
            cfg = "host",
        ),
        "_importer": attr.label(
            executable = True,
            default = Label("@importer//file:downloaded"),
            cfg = "host",
        ),
    },
    implementation = _impl,
)
