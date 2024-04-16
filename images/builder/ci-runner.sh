#!/usr/bin/env bash

# Copyright 2018 The cert-manager Authors.
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

BUILD_DIR="${1:-}"
if [ -z "${BUILD_DIR}" ]; then
    echo "Invalid usage. Use as $0 path/to/build/dir [additional arguments]"
    exit 1
fi
shift

WORKSPACE="$(bazel info workspace)"

echo "Activating service account..."
gcloud auth activate-service-account --key-file="${GOOGLE_APPLICATION_CREDENTIALS}"

echo "Generating docker credentials..."
gcloud auth configure-docker --quiet

echo "Executing builder..."
PUSHED_IMAGE=$(bazel run \
    //images/builder -- \
    --build-dir "${WORKSPACE}"/"${BUILD_DIR}" "$@")

echo "Build complete!"

if [ -z "${PUSHED_IMAGE}" ]; then
    echo "No image pushed to registry"
    exit 0
fi

echo "Pushed image ${PUSHED_IMAGE}"
echo

user="${GITHUB_USER:-}"
token="${GITHUB_TOKEN_FILE:-}"
if [ -z "${user}" ] || [ -z "${token}" ]; then
    echo "Skipping patching job configs"
    exit 0
fi

echo "Patching YAML files for new image"
find "${WORKSPACE}/config/jobs" -type f -name '*.yaml' | \
    xargs bazel run //tools/image-bumper -- \
    --image-regex "${PUSHED_IMAGE}"

ensure-config() {
  local username="jetstack-bot"
  local email="jetstack-bot@users.noreply.github.com"
  echo "git config user.name=$username user.email=$email..." >&2
  git config user.name "$username"
  git config user.email "$email"
}
ensure-config "$@"

image_name=$(basename "${PUSHED_IMAGE}")
title="Automatic bump of ${image_name} jobs"
git add -A
git commit -s -m "${title}"
git push -f "git@github.com:${user}/testing.git" HEAD:autobump-"${image_name}"

bazel run @test_infra//robots/pr-creator -- \
    --github-token-path="${token}" \
    --org jetstack --repo testing --branch master \
    --title="${title}" --match-title="Bump ${image_name} jobs" \
    --body="Automatically bumped jobs that referenced image \`${PUSHED_IMAGE}\`\n\n/kind bump" \
    --source="${user}":autobump-"${image_name}" \
    --confirm

echo "Complete!"
