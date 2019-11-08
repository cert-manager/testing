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

# Includes golang, docker-in-docker and gcloud
ARG BASE_IMAGE
FROM ${BASE_IMAGE}
LABEL maintainer="joshua.vanleeuwen@jetstack.io"

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
    && apt-get clean \
    && python -m pip install --upgrade pip setuptools wheel

# install golang
ARG GO_VERSION
RUN wget -O go${GO_VERSION}.amd64.tar.gz https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go${GO_VERSION}.amd64.tar.gz
ENV PATH "/usr/local/go/bin:${PATH}"
