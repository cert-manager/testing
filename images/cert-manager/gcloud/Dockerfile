# Copyright 2017 The Kubernetes Authors.
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

# Includes go and gcloud
FROM golang:1.9.3
LABEL maintainer="senlu@google.com"

# add env we can debug with the image name:tag
ARG IMAGE_ARG
ENV IMAGE=${IMAGE_ARG}

RUN apt-get update && apt-get install -y \
    python \
    rsync \
    wget && \
    apt-get clean

ENV GCLOUD_VERSION 163.0.0
RUN wget https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-$GCLOUD_VERSION-linux-x86_64.tar.gz && \
    tar xf google-cloud-sdk-$GCLOUD_VERSION-linux-x86_64.tar.gz -C / && \
    rm google-cloud-sdk-$GCLOUD_VERSION-linux-x86_64.tar.gz && \
    /google-cloud-sdk/install.sh && \
    curl -Lo helm.tar.gz https://storage.googleapis.com/kubernetes-helm/helm-v2.10.0-linux-amd64.tar.gz && \
    tar xvf helm.tar.gz && \
    mv linux-amd64/helm /usr/local/bin
ENV PATH "/google-cloud-sdk/bin:${PATH}"

RUN curl -L https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 >/tmp/dep \
    && curl -L https://github.com/mikefarah/yaml/releases/download/1.13.1/yaml_linux_amd64 >/tmp/yaml \
    && curl -L vert https://github.com/Masterminds/vert/releases/download/v0.1.0/vert-v0.1.0-linux-amd64 >/tmp/vert \
    && chmod +x /tmp/dep \
    && chmod +x /tmp/yaml \
    && chmod +x /tmp/vert \
    && mv /tmp/dep /tmp/yaml /tmp/vert /usr/local/bin/
RUN go get -v golang.org/x/tools/cmd/goimports

WORKDIR /workspace
ADD runner /
ENTRYPOINT ["/bin/bash", "/runner"]
