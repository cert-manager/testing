#!/usr/bin/env bash

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

PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../.." && pwd )"
SCRIPT_DIR="${PROJECT_DIR}/images/image-builder-script"

BUILD_DIR="${1:-}"
if [ -z "${BUILD_DIR}" ]; then
    echo "Invalid usage. Use as $0 path/to/build/dir [additional arguments]"
    exit 1
fi
shift

if [ -n "${GOOGLE_APPLICATION_CREDENTIALS:-}" ]; then
    echo "GOOGLE_APPLICATION_CREDENTIALS set, using service account"

    echo "Activating service account..."
    gcloud auth activate-service-account --key-file="${GOOGLE_APPLICATION_CREDENTIALS}"

    echo "Generating docker credentials..."
    gcloud auth configure-docker --quiet
else
    echo "WARNING: GOOGLE_APPLICATION_CREDENTIALS not set"
fi

echo "Executing builder..."
PUSHED_IMAGE=$(cd "$SCRIPT_DIR" && \
    go run . --build-dir "${PROJECT_DIR}"/"${BUILD_DIR}" $@)

echo "Build complete!"

if [ -z "${PUSHED_IMAGE}" ]; then
    echo "No image pushed to registry"
    exit 0
fi

echo "Pushed image ${PUSHED_IMAGE}"
echo
