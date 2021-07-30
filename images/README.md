# testing/images

Images used for various CI jobs for cert-manager and related projects.

All images are built in Prow. There is a Prow job per image in [config/jobs/testing/test-trusted.yaml](../config/jobs/testing/testing-trusted.yaml).

Most images are built using the scripts in [images/builder](./builder).

### kind

[images/kind](./kind) contains a script to build our own kind image—this may be useful when needing to test against a particular version of Kubernetes for which there is no 'official' kind image yet.

To build for a new Kubernetes version, change the `KUBERNETES_VERSION` variable in `images/kind/build.sh`—­this should trigger the `post-testing-push-kind` Prow job and a `gcr.io/jetstack-build-infra-images/kind:<KUBERNETES_VERSION>` image should get built and pushed.

This image build does not use [images/builder](./builder) functionality because the image is built with `kind`, not Docker.
