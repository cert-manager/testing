# Copyright 2021 The cert-manager Authors.
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

ARG BASE_IMAGE=golang:buster
FROM ${BASE_IMAGE}

ARG TERRAFORM_VERSION=1.0.0
ARG BAZEL_VERSION=4.0.0
ENV USE_BAZEL_VERSION=${BAZEL_VERSION}

RUN mkdir -p /tmp/setup && \
    cd /tmp/setup && \
    apt-get -y update && \
    apt-get install -y python3 unzip --no-install-recommends && \
    curl https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip -o awscliv2.zip && \
    unzip awscliv2.zip && \
    ./aws/install && \
    curl -o aws-iam-authenticator https://amazon-eks.s3.us-west-2.amazonaws.com/1.19.6/2021-01-05/bin/linux/amd64/aws-iam-authenticator && \
    chmod +x ./aws-iam-authenticator && \
    mv aws-iam-authenticator /usr/local/bin/ && \
    curl -O https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip &&  \
    mv terraform /usr/local/bin/ && \
    go get github.com/bazelbuild/bazelisk && \
    bazelisk && \
    ln -s /go/bin/bazelisk /usr/local/bin/bazel && \
    cd && \
    rm -rf /tmp/setup && \
    aws --version && \
    aws-iam-authenticator version && \
    bazel --version && \
    terraform version
