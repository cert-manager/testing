#!/usr/bin/env bash

# Copyright 2026 The cert-manager Authors.
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

set -euo pipefail

PROW_URL="https://prow.infra.cert-manager.io"
PROW_JOBS_URL="${PROW_URL}/prowjobs.js?omit=annotations,labels,decoration_config,pod_spec"

usage() {
    echo "Usage: $0 <list|rerun> <search-term>"
    echo
    echo "Commands:"
    echo "  list   List the latest run of each periodic job matching <search-term>"
    echo "  rerun  Rerun all non-success periodic jobs matching <search-term>"
    echo
    echo "Examples:"
    echo "  $0 list release-1.21"
    echo "  $0 rerun release-1.21"
    exit 1
}

if [[ $# -ne 2 ]]; then
    usage
fi

command="$1"
search_term="$2"

fetch_latest_periodics() {
    curl -s "${PROW_JOBS_URL}" \
        | jq -r --arg q "${search_term}" '
            [.items[] | select(.spec.type == "periodic" and (.spec.job | contains($q)))]
            | group_by(.spec.job)
            | map(sort_by(.status.startTime) | last)
        '
}

case "${command}" in
    list)
        fetch_latest_periodics | jq -r '
            sort_by(.spec.job)[]
            | "\(.spec.job): \(.status.state) (started: \(.status.startTime), id: \(.metadata.name))"
        '
        ;;
    rerun)
        fetch_latest_periodics | jq -r '
            map(select(.status.state != "success"))[]
            | "\(.spec.job)\t\(.metadata.name)"
        ' \
        | while IFS=$'\t' read -r job id; do
            echo "Rerunning: ${job}"
            kubectl create -f "${PROW_URL}/rerun?mode=original&prowjob=${id}"
        done
        ;;
    *)
        usage
        ;;
esac
