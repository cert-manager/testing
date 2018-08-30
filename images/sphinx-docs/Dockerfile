FROM debian:jessie
LABEL maintainer="christian.simon@jetstack.io"

RUN apt-get update && apt-get install -y \
    git \
    curl \
    build-essential \
    python-minimal \
    python-pip \
    python3-enchant \
    wamerican \
    wbritish && \
    apt-get clean

ENV GCLOUD_VERSION 178.0.0
ENV GCLOUD_HASH 2e0bbbf81c11164bf892cf0b891751ba4e5172661eff907ad1f7fc0b6907e296
RUN curl -sL -o /tmp/google-cloud-sdk.tar.gz https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-$GCLOUD_VERSION-linux-x86_64.tar.gz && \
    echo "$GCLOUD_HASH  /tmp/google-cloud-sdk.tar.gz" | sha256sum -c && \
    tar xf /tmp/google-cloud-sdk.tar.gz -C / && \
    rm /tmp/google-cloud-sdk.tar.gz && \
    /google-cloud-sdk/install.sh
ENV PATH "/google-cloud-sdk/bin:${PATH}"

RUN pip install virtualenv

ADD runner /

# run as non root
ENV VENV_PATH /venv
RUN mkdir -p /home/nobody /workspace ${VENV_PATH} && chown nobody /home/nobody /workspace ${VENV_PATH}

WORKDIR /workspace

ENV HOME /home/nobody
USER nobody

ADD requirements.txt /tmp/requirements.txt
RUN virtualenv -p $(which python3) ${VENV_PATH}
RUN ${VENV_PATH}/bin/pip install -r /tmp/requirements.txt

ENTRYPOINT ["/bin/bash", "/runner"]
