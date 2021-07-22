#!/usr/bin/env bash

# Copyright 2021 The Jetstack contributors.
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

set -o errexit
set -o nounset
set -o pipefail

# Tag to check out in k/k repo. Kind will build Kubernetes binaries from that
# tag and include in the built KIND image.
KUBERNETES_VERSION=v1.22.0-beta.2
# Version of the kind CLI to use to build the kind image.
KIND_BASE_VERSION=v0.11.1

echo "Downloading dependencies..."

go get sigs.k8s.io/kind@${KIND_BASE_VERSION}
export PATH=$(go env GOPATH)/bin:$PATH

# go get seems to not work for k/k see https://github.com/kubernetes/kubernetes/issues/79384
kube_path=$(go env GOPATH)/src/k8s.io/kubernetes
mkdir -p $kube_path
git clone --branch ${KUBERNETES_VERSION} \
	--depth 1 \
	https://github.com/kubernetes/kubernetes \
	${kube_path}

image_tag=gcr.io/jetstack-build-infra-images/kind:${KUBERNETES_VERSION}

echo "Building $image_tag..."
kind build node-image \
	--image ${image_tag}

echo "Activating service account..."
gcloud auth activate-service-account --key-file="${GOOGLE_APPLICATION_CREDENTIALS}"

echo "Generating docker credentials..."
gcloud auth configure-docker --quiet

echo "Pushing ${image_tag}..."
docker push ${image_tag}

echo "${image_tag} built and pushed!"
