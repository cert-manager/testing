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

FROM debian:buster-slim

RUN apt-get -y update \
&& apt-get -y install --no-install-recommends \
   bash \
   ca-certificates \
   curl \
   file \
   findutils \
   gawk \
   git \
   gnupg \
   jq \
   python3 \
   python3-yaml \
   util-linux \
   uuid-runtime \
&& curl -fsSL https://deb.nodesource.com/setup_8.x | bash - \
&& apt-get -y install --no-install-recommends nodejs \
&& npm i markdown-spellcheck -g \
&& curl -fsSL https://www.katacoda.com/cli/download/linux -o /usr/local/bin/katacoda \
&& chmod +x /usr/local/bin/katacoda \
&& rm -rf /var/lib/apt/lists/*
