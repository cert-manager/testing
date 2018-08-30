# Prow deployment

This directory contains manifests used for the deployment of the Prow cluster.

After modifying files here, run `bazel run //prow/cluster:production.apply` to
apply changes to the production build cluster.

## TODO

Instead of explicitly specifying image tags in the deployment manifests, we can
instead automatically build and push docker images based on the chosen revision
of jetstack/test-infra we have defined in the WORKSPACE file.

This will make it easier to update Prow, and also mean we can easily control what
version of test-infra is deployed to our production cluster.

We can use [rules_k8s](https://github.com/bazelbuild/rules_k8s) for this.
