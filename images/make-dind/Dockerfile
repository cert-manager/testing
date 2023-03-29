# Copyright 2023 The Jetstack contributors.
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

# Includes make, docker-in-docker and gcloud
ARG DEBIAN_VERSION
FROM debian:"${DEBIAN_VERSION}"

LABEL maintainer="cert-manager-maintainers@googlegroups.com"

#
# BEGIN: DOCKER IN DOCKER SETUP
#

# Install Docker deps, some of these are already installed in the image but
# that's fine since they won't re-install and we can reuse the code below
# for another image someday.
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        apt-transport-https \
        ca-certificates \
        curl \
        gnupg2 \
        software-properties-common \
        lsb-release \
    && apt-get clean

# Add the Docker apt-repository
RUN mkdir -p /etc/apt/keyrings \
    && curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg \
    && echo "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" \
    | tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker
# TODO(bentheelder): the `sed` is a bit of a hack, look into alternatives.
# Why this exists: `docker service start` on debian runs a `cgroupfs_mount` method,
# We're already inside docker though so we can be sure these are already mounted.
# Trying to remount these makes for a very noisy error block in the beginning of
# the pod logs, so we just comment out the call to it... :shrug:
ARG DOCKER_VERSION
RUN apt-get update \
    && apt-get install -y --no-install-recommends docker-ce="${DOCKER_VERSION}" \
    && apt-get clean \
    && sed -i 's/cgroupfs_mount$/#cgroupfs_mount\n/' /etc/init.d/docker \
    && update-alternatives --set iptables /usr/sbin/iptables-legacy \
    && update-alternatives --set ip6tables /usr/sbin/ip6tables-legacy

# Move Docker's storage location & enable experimental features & add Container Registry cache (see: https://cloud.google.com/container-registry/docs/pulling-cached-images)
RUN echo 'DOCKER_OPTS="${DOCKER_OPTS} --data-root=/docker-graph --experimental --registry-mirror=https://mirror.gcr.io"' | \
    tee --append /etc/default/docker

# NOTE this should be mounted and persisted as a volume ideally (!)
VOLUME /docker-graph

#
# END: DOCKER IN DOCKER SETUP
#

# Add the google-cloud-sdk apt-repository
RUN mkdir -p /etc/apt/keyrings \
    && curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | gpg --dearmor -o /etc/apt/keyrings/cloud.google.gpg \
    && echo "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt "cloud-sdk-$(. /etc/os-release && echo "$VERSION_CODENAME")" main" \
    | tee /etc/apt/sources.list.d/google-cloud-sdk.list > /dev/null

# moreutils is used to get timestamping on stdout
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        google-cloud-sdk \
        zip \
        unzip \
        python \
        python3-pip \
        wget \
        ca-certificates \
        git \
        make \
        rsync \
        patch \
    && apt-get clean \
    && python3 -m pip install --upgrade pip setuptools wheel

WORKDIR /workspace

COPY runner /usr/local/bin/runner

ENTRYPOINT ["/usr/local/bin/runner"]
