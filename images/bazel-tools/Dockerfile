# Copyright 2021 The Jetstack contributors.
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

# Includes bazel, make, node, go, goversion, jq
ARG BASE_IMAGE

FROM ${BASE_IMAGE}

LABEL maintainer="cert-manager-maintainers@googlegroups.com"

COPY semversort.sh /usr/local/bin/semversort

ARG NODE_VERSION
# install goversion, gcrane, cosign, gh cli, jq and node
RUN go install github.com/cert-manager/goversion@v1.3.0 && \
    go install github.com/sigstore/cosign/cmd/cosign@v1.10.1 && \
    go install github.com/google/go-containerregistry/cmd/gcrane@v0.11.0 && \
    apt-get update && \
    apt-get install -y \
    nodejs=${NODE_VERSION} && \
    curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list > /dev/null && \
    apt update && \
    apt install gh=2.14.4

# Add GOPATH/bin to PATH
ENV PATH=/root/go/bin:$PATH
