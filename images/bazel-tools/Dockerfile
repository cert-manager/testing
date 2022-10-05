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
# install goversion, gcrane, cosign, gh cli, jq, yq, helm and node
RUN go install github.com/cert-manager/goversion@v1.3.0 && \
    go install github.com/sigstore/cosign/cmd/cosign@v1.10.1 && \
    go install github.com/google/go-containerregistry/cmd/gcrane@v0.11.0 && \
    apt-get update && \
    apt-get install -y \
    nodejs=${NODE_VERSION} && \
    curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg && \
    curl -O -L https://get.helm.sh/helm-v3.10.0-linux-amd64.tar.gz && \
    tar -xvzf helm-v3.10.0-linux-amd64.tar.gz && \
    mv linux-amd64/helm /usr/local/bin/ && \
    curl -O https://github.com/mikefarah/yq/releases/download/v4.27.3/yq_linux_amd64 && \
    chmod +x yq_linux_amd64 && \
    mv yq_linux_amd64 /usr/local/bin/yq && \
    curl -O -L https://github.com/cli/cli/releases/download/v2.17.0/gh_2.17.0_linux_amd64.tar.gz && \
    tar -xvzf gh_2.17.0_linux_amd64.tar.gz && \
    mv gh_2.17.0_linux_amd64/bin/gh /usr/local/bin/

# Add GOPATH/bin to PATH
ENV PATH=/root/go/bin:$PATH
