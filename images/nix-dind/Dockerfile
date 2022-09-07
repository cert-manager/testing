# Copyright 2022 The Jetstack contributors.
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

# Includes docker-in-docker and gcloud
ARG BASE_IMAGE
FROM ${BASE_IMAGE}
LABEL maintainer="cert-manager-maintainers@googlegroups.com"

# xz is required for nix curl install.
RUN apt-get install -y xz-utils

# install nix.
ARG NIX_VERSION
RUN mkdir -m 0755 /nix && groupadd -r nixbld && chown root /nix
RUN for n in $(seq 1 10); do useradd -c "Nix build user $n" -d /var/empty -g nixbld -G nixbld -M -N -r -s "$(command -v nologin)" "nixbld$n"; done
RUN mkdir -p /etc/nix && echo "experimental-features = nix-command flakes\nsandbox = false" > /etc/nix/nix.conf
RUN bash -c "$(curl -L https://releases.nixos.org/nix/nix-${NIX_VERSION}/install) --no-daemon"

# add nix store to path.
ENV PATH=/root/.nix-profile/bin:$PATH
