
#!/usr/bin/env bash
# Copyright 2021 The Jetstack contributors.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# */

# Runs prow/pj-on-kind.sh with config arguments specific to Jetstack Prow config.
# Requries go, docker, and kubectl.

# Copied and adapted from https://github.com/istio/test-infra/blob/master/prow/pj-on-kind.sh

# Example usage:
# ./prow/pj-on-kind.sh ci-cert-manager-e2e-v1-21

set -o errexit
set -o nounset
set -o pipefail

set -x
SCRIPT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"
export REPO_ROOT="${SCRIPT_ROOT}/.."

export CONFIG_PATH="${REPO_ROOT}/config/config.yaml"
export JOB_CONFIG_PATH="${REPO_ROOT}/config/jobs"

bash <(curl -sSfL https://raw.githubusercontent.com/kubernetes/test-infra/master/prow/pj-on-kind.sh) "$@"