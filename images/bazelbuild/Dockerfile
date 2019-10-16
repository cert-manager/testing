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

# Includes bazel, docker-in-docker and gcloud
ARG DEBIAN_VERSION=stretch
FROM debian:"${DEBIAN_VERSION}"
LABEL maintainer="james@jetstack.io"

#
# BEGIN: DOCKER IN DOCKER SETUP
#

# Install Docker deps, some of these are already installed in the image but
# that's fine since they won't re-install and we can reuse the code below
# for another image someday.
RUN apt-get update && apt-get install -y --no-install-recommends \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg2 \
    software-properties-common \
    lsb-release

# Add the Docker apt-repository
RUN curl -fsSL https://download.docker.com/linux/$(. /etc/os-release; echo "$ID")/gpg \
    | apt-key add - && \
    add-apt-repository \
    "deb [arch=amd64] https://download.docker.com/linux/$(. /etc/os-release; echo "$ID") \
    $(lsb_release -cs) stable"

# Install Docker
# TODO(bentheelder): the `sed` is a bit of a hack, look into alternatives.
# Why this exists: `docker service start` on debian runs a `cgroupfs_mount` method,
# We're already inside docker though so we can be sure these are already mounted.
# Trying to remount these makes for a very noisy error block in the beginning of
# the pod logs, so we just comment out the call to it... :shrug:
# TODO(benthelder): update docker version. This is pinned because of
# https://github.com/kubernetes/test-infra/issues/6187
ARG DOCKER_VERSION="17.09.1~ce-0~debian"
RUN apt-get update && \
    apt-get install -y --no-install-recommends docker-ce="${DOCKER_VERSION}" && \
    sed -i 's/cgroupfs_mount$/#cgroupfs_mount\n/' /etc/init.d/docker

# Move Docker's storage location & enable experimental features
RUN echo 'DOCKER_OPTS="${DOCKER_OPTS} --data-root=/docker-graph --experimental"' | \
    tee --append /etc/default/docker
# NOTE this should be mounted and persisted as a volume ideally (!)
# We will make a fallback one now just in case
RUN mkdir /docker-graph

#
# END: DOCKER IN DOCKER SETUP
#

# Add new repos to install google-cloud-sdk
RUN echo "deb http://packages.cloud.google.com/apt cloud-sdk-$(lsb_release -c -s) main" | \
    tee -a /etc/apt/sources.list.d/google-cloud-sdk.list

RUN curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -

# make is installed simply because a lot of things use it - it is not required
# by Bazel
RUN apt-get update && apt-get install -y --no-install-recommends \
    google-cloud-sdk \
    pkg-config \
    zip \
    g++ \
    zlib1g-dev \
    unzip \
    python \
    python-pip \
    wget \
    ca-certificates \
    git \
    mercurial \
    make \
    rsync \
    patch \
    openssh-client \
    && apt-get clean \
    && python -m pip install --upgrade pip setuptools wheel

ARG BAZEL_VERSION
ARG BAZEL_CHANNEL=release
ARG INSTALLER_VERSION_NAME="${BAZEL_VERSION}"

RUN INSTALLER="bazel-${INSTALLER_VERSION_NAME}-installer-linux-x86_64.sh"; \
    DOWNLOAD_URL="https://releases.bazel.build/${BAZEL_VERSION}/${BAZEL_CHANNEL}/${INSTALLER}" \
    wget -q "${DOWNLOAD_URL}" && \
    chmod +x "${INSTALLER}" && "./${INSTALLER}" && rm "${INSTALLER}"

COPY runner /usr/local/bin/runner
COPY create_bazel_cache_rcs.sh /usr/local/bin/create_bazel_cache_rcs.sh
COPY coalesce.py /usr/local/bin/coalesce.py

WORKDIR /workspace
