#!/bin/bash

# Change USERNAME to your username (e.g. client or crierclient)
USERNAME=crierclient

CSR_FILE=$USERNAME.csr
KEY_FILE=$USERNAME.key
CRT_FILE=$USERNAME.crt

openssl genrsa -out $KEY_FILE 2048
openssl req -new -key $KEY_FILE -out $CSR_FILE -subj "/CN=$USERNAME"

cat <<EOF | kubectl create -f -
apiVersion: certificates.k8s.io/v1
kind: CertificateSigningRequest
metadata:
    name: $USERNAME
spec:
    signerName: kubernetes.io/kube-apiserver-client
    groups:
    - system:authenticated
    request: $(cat $CSR_FILE | base64 | tr -d '\n')
    usages:
    - digital signature
    - key encipherment
    - client auth
EOF

kubectl certificate approve $USERNAME
kubectl get csr $USERNAME -o jsonpath='{.status.certificate}' | base64 -d > $CRT_FILE
