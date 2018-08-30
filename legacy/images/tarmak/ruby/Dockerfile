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

ARG RUBY_VERSION
FROM ruby:${RUBY_VERSION}
LABEL maintainer="christian.simon@jetstack.io"

RUN apt-get update && apt-get install -y \
    python \
    git \
    curl \
    zlib1g-dev \
    build-essential \
    libc++-dev \
    ruby-dev \
    libssl-dev \
    libreadline-dev \
    libyaml-dev \
    libsqlite3-dev \
    sqlite3 \
    libxml2-dev \
    libxslt1-dev \
    libvirt-dev \
    libcurl4-openssl-dev \
    python3-software-properties \
    libffi-dev \
    nodejs \
    rsync \
    wget && \
    apt-get clean

# Download and install vagrant
ENV VAGRANT_VERSION 2.1.1
ENV VAGRANT_HASH cee18b6f3b1209ea5878c22cfd84a9f0004f20ef146cb7a18aada19162928a0f
RUN curl -sL -o /tmp/vagrant.deb https://releases.hashicorp.com/vagrant/${VAGRANT_VERSION}/vagrant_${VAGRANT_VERSION}_x86_64.deb \
 && echo "${VAGRANT_HASH}  /tmp/vagrant.deb" | sha256sum -c \
 && DEBIAN_FRONTEND=noninteractive dpkg -i /tmp/vagrant.deb \
 && rm /tmp/vagrant.deb

# Download and install vagrant's libvirt plugin
RUN vagrant plugin install vagrant-libvirt

ENV GCLOUD_VERSION 206.0.0
ENV GCLOUD_HASH d39293914b2e969bfe18dd19eb77ba96d283995f8cf1e5d7ba6ac712a3c9479a
RUN curl -sL -o /tmp/google-cloud-sdk.tar.gz https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-$GCLOUD_VERSION-linux-x86_64.tar.gz && \
    echo "$GCLOUD_HASH  /tmp/google-cloud-sdk.tar.gz" | sha256sum -c && \
    tar xf /tmp/google-cloud-sdk.tar.gz -C / && \
    rm /tmp/google-cloud-sdk.tar.gz && \
    /google-cloud-sdk/install.sh
ENV PATH "/google-cloud-sdk/bin:${PATH}"

WORKDIR /workspace
ADD runner /

# add env we can debug with the image name:tag
ARG IMAGE_ARG
ENV IMAGE=${IMAGE_ARG}

ENTRYPOINT ["/bin/bash", "/runner"]
