#!/bin/bash

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

set -o errexit
set -o nounset
set -o pipefail

# This script should be executed via bazel only.

OUT=$1
shift
DOCKERFILE=$1
shift
ROOTDIR=$(dirname "${DOCKERFILE}")

# TODO: construct a directory to contain the dockerfile's build context
# This will need to be assembled based on paths passed to this script.
# Without doing this, we cannot automatically rebuild when changes are
# made to files that are added into images, as we currently don't explicitly
# state which workspace files are required.

# TODO: randomise name here
docker build -t "testimagebuild" "$@" -f "$DOCKERFILE" "$ROOTDIR"
docker save "testimagebuild" -o "$OUT"
