#!/bin/bash

# Copyright 2023 The Jetstack contributors.
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
