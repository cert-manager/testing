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

IMG = eu.gcr.io/jetstack-build-infra-images/tarmak-ruby
TAG := $(shell date +v%Y%m%d)-$(shell git describe --tags --always --dirty)

VERSION = 2.4.4

all: build

do_build = \
	docker build --pull --build-arg RUBY_VERSION=$(VERSION) --build-arg IMAGE_ARG=$(IMG):$(TAG)-$(VERSION) -t $(IMG):$(TAG)-$(VERSION) .;\
	docker tag $(IMG):$(TAG)-$(VERSION) $(IMG):latest-$(VERSION);\
	echo Built $(IMG):$(TAG)-$(VERSION) and tagged with latest-$(VERSION);

do_push = \
	docker push $(IMG):$(TAG)-$(VERSION);\
	docker push $(IMG):latest-$(VERSION);\
	echo Pushed $(IMG) with $(TAG)-$(VERSION) and latest-$(VERSION) tags;

build:
	@$(foreach VERSION,$(VERSION),$(do_build))

push: build
	@$(foreach VERSION,$(VERSION),$(do_push))

.PHONY: all build push
