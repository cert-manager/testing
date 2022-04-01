FROM debian

# Usage:
#
#   nerdctl build -t local/make-dind -f images/make-dind/Dockerfile ./images/make-dind
#   nerdctl run -it --privileged -v /tmp:/tmp --tmpfs /var/lib/containerd local/make-dind bash

RUN apt update \
    && apt install curl bash make curl python3 perl jq git tar grep iptables -y \
    && curl -sSL https://github.com/containerd/nerdctl/releases/download/v0.17.1/nerdctl-full-0.17.1-linux-amd64.tar.gz | tar zx -C / \
    && mkdir -p /opt/cni/bin \
    && curl -sSL https://github.com/containernetworking/plugins/releases/download/v1.1.1/cni-plugins-linux-amd64-v1.1.1.tgz | tar xz -C /opt/cni/bin
COPY docker-entrypoint.sh /
VOLUME /var/lib/containerd
ENTRYPOINT ["/docker-entrypoint.sh"]
