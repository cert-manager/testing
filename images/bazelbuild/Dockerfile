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

# Includes bazel only
FROM ubuntu:16.04
LABEL maintainer="james@jetstack.io"

# add env we can debug with the image name:tag
ARG IMAGE_ARG
ENV IMAGE=${IMAGE_ARG}

ARG BAZEL_VERSION

RUN apt-get update && apt-get install -y --no-install-recommends \
    pkg-config \
    zip \
    g++ \
    zlib1g-dev \
    unzip \
    python \
    wget

RUN INSTALLER="bazel-${BAZEL_VERSION}-installer-linux-x86_64.sh"; \
    DOWNLOAD_URL="https://github.com/bazelbuild/bazel/releases/download/${BAZEL_VERSION}/${INSTALLER}"; \
    wget -q "${DOWNLOAD_URL}" && \
    chmod +x "${INSTALLER}" && "./${INSTALLER}" && rm "${INSTALLER}"

WORKDIR /workspace
