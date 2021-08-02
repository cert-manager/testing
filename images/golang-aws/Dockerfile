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
    terraform version && \
    bazel --version
