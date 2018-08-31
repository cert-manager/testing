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

load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_bundle")

def bootstrap_image(name, base, runner = "//legacy/bootstrap:runner"):
    container_image(
        name = name,
        base = base,
        directory = "/workspace",
        # TODO: lookup name of 'runner' target
        entrypoint = ["/bin/bash", "/runner"],
        files = [runner],
    )

def bootstrap_image_bundle(name, images, stamp):
    i = 0
    new_images = {}
    for k, v in images.items():
        new_name = "%s.%d" % (name, i)
        bootstrap_image(new_name, v)
        new_images[k] = ":%s" % new_name
        i += 1

    container_bundle(
        name = name,
        images = new_images,
        stamp = stamp,
    )
