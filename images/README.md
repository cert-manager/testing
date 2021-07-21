# testing/images

Images used for various CI jobs for cert-manager and related projects.

All images are built in Prow. There is a Prow job per image in [config/jobs/testing/test-trusted.yaml](config/jobs/testing/test-trusted.yaml).

Most images are built using the scripts in [image/builder](images/builder).

### kind

[images/kind](images/kind) contains a script to build our own KIND image- this may be useful when needing to test against a particular version of Kubernetes for which there is no 'official' KIND image yet.

To build for a new Kubernetes version, change the  in images/kind/build.sh - this should trigger `post-testing-push-kind` Prow job and a `gcr.io/jetstack-build-infra-images/kind:<KUBERNETES_VERSION>` image should get built and pushed.

This image build does not use [image/builder](images/builder) functionality because the image is built with KIND CLI, not Docker.