FROM golang:1.16-buster

RUN apt-get -y update \
    && apt-get install -y python3 unzip --no-install-recommends; \
    curl https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip -o awscliv2.zip; \
    unzip awscliv2.zip; \
    ./aws/install; \
    curl -o aws-iam-authenticator https://amazon-eks.s3.us-west-2.amazonaws.com/1.19.6/2021-01-05/bin/linux/amd64/aws-iam-authenticator; \
    chmod +x ./aws-iam-authenticator; \
    mkdir -p $HOME/bin && cp ./aws-iam-authenticator $HOME/bin/aws-iam-authenticator && export PATH=$PATH:$HOME/bin; \
    curl -O https://releases.hashicorp.com/terraform/1.0.3/terraform_1.0.3_linux_amd64.zip; \
    unzip terraform_1.0.3_linux_amd64.zip;  \
    mv terraform /usr/local/bin/;