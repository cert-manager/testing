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

IMG = eu.gcr.io/jetstack-build-infra-images/sphinx-docs
TAG := $(shell date +v%Y%m%d)-$(shell git describe --tags --always --dirty)

# TODO: add sphinx version to image tag

all: build

do_build = \
	docker build --pull --build-arg IMAGE_ARG=$(IMG):$(TAG) -t $(IMG):$(TAG) .;\
	docker tag $(IMG):$(TAG) $(IMG):latest;\
	echo Built $(IMG):$(TAG) and tagged with latest;

do_push = \
	docker push $(IMG):$(TAG);\
	docker push $(IMG):latest;\
	echo Pushed $(IMG) with $(TAG) and latest tags;

build:
	@$(do_build)

push: build
	@$(do_push)

.PHONY: all build push
