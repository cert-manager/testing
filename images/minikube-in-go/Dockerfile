# Copyright 2017 Jetstack Ltd.

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

# Includes go, gcloud, minikube, kubectl and minikube
FROM debian:jessie
LABEL maintainer="james@jetstack.io"

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update \
    && apt-get install -y \
        libvirt-bin \
        qemu-kvm \
        libvirt-dev \
        curl \
        git \
        gcc \
        build-essential \
        ca-certificates \
        curl \
        file \
        git \
        jq \
        openssh-client \
        pkg-config \
        python \
        python-openssl \
        python-pip \
        rsync \
        unzip \
        wget \
        xz-utils \
        zip \
        zlib1g-dev \
    && apt-get clean

RUN curl -L https://github.com/docker/machine/releases/download/v0.13.0/docker-machine-Linux-x86_64 >/tmp/docker-machine \
    && chmod +x /tmp/docker-machine \
    && cp /tmp/docker-machine /usr/local/bin/docker-machine

RUN wget https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz && \
	tar -C /usr/local -xzf go1.9.1.linux-amd64.tar.gz
ENV PATH "/usr/local/go/bin:/go/bin:${PATH}"
ENV GOPATH "/go"

RUN go get -v github.com/dhiltgen/docker-machine-kvm/cmd/docker-machine-driver-kvm
RUN go get -v golang.org/x/tools/cmd/goimports

ENV GCLOUD_VERSION 163.0.0
RUN wget https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-$GCLOUD_VERSION-linux-x86_64.tar.gz && \
    tar xf google-cloud-sdk-$GCLOUD_VERSION-linux-x86_64.tar.gz -C / && \
    rm google-cloud-sdk-$GCLOUD_VERSION-linux-x86_64.tar.gz && \
    /google-cloud-sdk/install.sh
ENV PATH "/google-cloud-sdk/bin:${PATH}"

ARG KUBERNETES_VERSION
ENV KUBERNETES_VERSION=$KUBERNETES_VERSION
RUN curl -Lo helm.tar.gz https://storage.googleapis.com/kubernetes-helm/helm-v2.10.0-linux-amd64.tar.gz && tar xvf helm.tar.gz && mv linux-amd64/helm /usr/local/bin \
    && curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/${KUBERNETES_VERSION}/bin/linux/amd64/kubectl && chmod +x kubectl && mv kubectl /usr/local/bin/ \
    && curl -Lo minikube https://storage.googleapis.com/minikube/releases/v0.23.0/minikube-linux-amd64 && chmod +x minikube && mv minikube /usr/local/bin/
ENV PATH "/usr/local/bin:${PATH}"

RUN curl -L https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 >/tmp/dep \
    && chmod +x /tmp/dep \
    && mv /tmp/dep /usr/local/bin/dep

# Install docker
# Note: 1.11+ changes the tarball format
RUN curl -L "https://get.docker.com/builds/Linux/x86_64/docker-1.9.1.tgz" \
    | tar -C /usr/bin -xvzf- --strip-components=3 usr/local/bin/docker

ARG BAZEL_VERSION

RUN INSTALLER="bazel-${BAZEL_VERSION}-installer-linux-x86_64.sh"; \
    DOWNLOAD_URL="https://github.com/bazelbuild/bazel/releases/download/${BAZEL_VERSION}/${INSTALLER}"; \
    wget -q "${DOWNLOAD_URL}" && \
    chmod +x "${INSTALLER}" && "./${INSTALLER}" && rm "${INSTALLER}"

WORKDIR /workspace
ADD runner /
ENTRYPOINT ["/bin/bash", "/runner"]
