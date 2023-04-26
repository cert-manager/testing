This folder contains the script that can be used to generate & sign client certificates used by the Prow jobs to authenticate against the Kubernetes API server.

Use the following command to update the kubeconfig file (https://kubernetes.io/docs/setup/best-practices/certificates/#configure-certificates-for-user-accounts):

```bash
KUBECONFIG=./crier-kubeconfig.yaml kubectl config set-credentials trusted --client-key crierclient.key --client-certificate crierclient.crt --embed-certs
```

Use the following commands to update the kubeconfig secrets in the clusters:

```bash
kubectl apply --server-side secret kubeconfig --from-file=config=kubeconfig.yaml
kubectl apply --server-side secret crier-kubeconfig --from-file=config=crier-kubeconfig.yaml

kubectl create secret generic crier-kubeconfig --from-file=config=crier-kubeconfig.yaml --dry-run=client -o yaml | kubectl apply --server-side -f -
```
