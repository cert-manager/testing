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

# Includes make and docker-in-docker

ARG DEBIAN_VERSION
FROM debian:"${DEBIAN_VERSION}"

LABEL maintainer="cert-manager-maintainers@googlegroups.com"

# Some of these deps might already be installed in the base image but we
# ensure they're installed here to ensure consistency
# TODO(SgtCoDFish): The python dependency can be removed once we remove the use of Python
# for scanning boilerplate. That saves about 50MB from the final image size.
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        apt-transport-https \
        ca-certificates \
        curl \
        gnupg2 \
        software-properties-common \
        lsb-release \
        g++ \
        zip \
        unzip \
        python3 \
        python3-pip \
        python3-setuptools \
        python3-wheel \
        wget \
        git \
        make \
        rsync \
        patch \
        jq \
        rsync \
    && apt-get clean

#
# BEGIN: DOCKER IN DOCKER SETUP
#

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
    && apt-get install -y --no-install-recommends \
        docker-ce="${DOCKER_VERSION}" \
        docker-buildx-plugin \
    && apt-get clean \
    && sed -i 's/cgroupfs_mount$/#cgroupfs_mount\n/' /etc/init.d/docker \
    && update-alternatives --set iptables /usr/sbin/iptables-legacy \
    && update-alternatives --set ip6tables /usr/sbin/ip6tables-legacy

# Move Docker's storage location & add Container Registry cache (see: https://cloud.google.com/container-registry/docs/pulling-cached-images)
# @inteon: added --mtu 1460 to fix network issues due to parent mtu < child mtu (see https://blog.zespre.com/dind-mtu-size-matters.html)
#          at the time of writing, the parent mtu is 1460 (see https://cloud.google.com/kubernetes-engine/docs/concepts/network-overview)
RUN echo 'DOCKER_OPTS="${DOCKER_OPTS} --mtu=1460 --data-root=/docker-graph --registry-mirror=https://mirror.gcr.io"' | \
    tee --append /etc/default/docker

# NOTE this should be mounted and persisted as a volume ideally (!)
VOLUME /docker-graph

#
# END: DOCKER IN DOCKER SETUP
#

WORKDIR /workspace

COPY runner /usr/local/bin/runner

ENTRYPOINT ["/usr/local/bin/runner"]
