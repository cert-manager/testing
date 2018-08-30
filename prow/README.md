# Prow deployment

This directory contains manifests used for the deployment of the Prow cluster.

## Updating Prow

If you make any changes to the manifests here, you will need to **manually**
run Bazel to apply these to production.

We do *not* currently automate the roll-out of the manifests in this repository.

You can 'apply' the latest changes in this repository using:

```
$ bazel run //prow/cluster:production.apply
```

In order to find connections details for the 'build-infra' and 'libvirt' clusters,
you will need to ensure you have two contexts already correctly configured.
Namely, `build-infra` and `libvirt`. Bazel will use the contexts with these names
to apply changes to Prow.

You can see where these context names are hardcoded in the [hack/print-workspace-status.sh](hack/print-workspace-status.sh)
file.

## TODO

Instead of explicitly specifying image tags in the deployment manifests, we can
instead automatically build and push docker images based on the chosen revision
of jetstack/test-infra we have defined in the WORKSPACE file.

This will make it easier to update Prow, and also mean we can easily control what
version of test-infra is deployed to our production cluster.

We can use [rules_k8s](https://github.com/bazelbuild/rules_k8s) for this.
