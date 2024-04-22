# Copyright 2023 The cert-manager Authors.
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

ARG BUILDER_IMAGE
ARG BASE_IMAGE

FROM ${BUILDER_IMAGE} as builder

WORKDIR /workspace

RUN git clone --depth=1 --branch option_recreate_prowjob_on_termination https://github.com/inteon/prow.git
RUN cd ./prow/ && CGO_ENABLED=0 go build -o ../prow-controller-manager ./cmd/prow-controller-manager/.

FROM ${BASE_IMAGE}
LABEL maintainer="cert-manager-maintainers@googlegroups.com"

WORKDIR /
COPY --from=builder /workspace/prow-controller-manager /prow-controller-manager

USER 65532:65532

ENTRYPOINT ["/prow-controller-manager"]
