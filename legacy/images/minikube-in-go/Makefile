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

IMG = eu.gcr.io/jetstack-build-infra-images/minikube-in-go
TAG := $(shell date +v%Y%m%d)-$(shell git describe --tags --always --dirty)

KUBERNETES_VERSIONS = v1.8.10 v1.7.15 v1.9.6
BAZEL_VERSION = 0.16.1

all: build

do_build = \
	docker build --pull --build-arg KUBERNETES_VERSION=$(VERSION) --build-arg BAZEL_VERSION=$(BAZEL_VERSION) --build-arg IMAGE_ARG=$(IMG):$(TAG)-$(VERSION) -t $(IMG):$(TAG)-$(VERSION) .;\
	docker tag $(IMG):$(TAG)-$(VERSION) $(IMG):latest-$(VERSION);\
	echo Built $(IMG):$(TAG)-$(VERSION) and tagged with latest-$(VERSION);

do_push = \
	docker push $(IMG):$(TAG)-$(VERSION);\
	docker push $(IMG):latest-$(VERSION);\
	echo Pushed $(IMG) with $(TAG)-$(VERSION) and latest-$(VERSION) tags;

build:
	$(foreach VERSION,$(KUBERNETES_VERSIONS),$(do_build))

push: build
	@$(foreach VERSION,$(KUBERNETES_VERSIONS),$(do_push))

.PHONY: all build push
