# Copyright 2018 The cert-manager Authors.
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

ARG BASE_IMAGE
FROM ${BASE_IMAGE}

# install golang
ARG GO_VERSION
RUN wget -O go${GO_VERSION}.amd64.tar.gz https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go${GO_VERSION}.amd64.tar.gz
ENV PATH "/usr/local/go/bin:${PATH}"
