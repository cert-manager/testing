FROM debian

# Usage:
#
#   nerdctl build -t local/make-dind -f images/make-dind/Dockerfile ./images/make-dind
#   nerdctl run -it --privileged -v /tmp:/tmp --tmpfs /var/lib/containerd local/make-dind bash

RUN apt update \
    && apt install curl bash make curl python3 perl jq git tar grep iptables -y \
    && rm -rf /var/lib/apt/lists/* \
    && curl -sSL https://github.com/containerd/nerdctl/releases/download/v1.2.1/nerdctl-full-1.2.1-linux-amd64.tar.gz | tar zx -C / \
    && ln -s /bin/nerdctl /bin/docker \
    && mkdir -p /opt/cni/bin \
    && curl -sSL https://github.com/containernetworking/plugins/releases/download/v1.2.0/cni-plugins-linux-amd64-v1.2.0.tgz | tar xz -C /opt/cni/bin

COPY docker-entrypoint.sh /
COPY containerd-cni.conflist.json /etc/cni/net.d/containerd-cni.conflist
COPY containerd.toml /etc/containerd/config.toml
VOLUME /var/lib/containerd
ENTRYPOINT ["/docker-entrypoint.sh"]
