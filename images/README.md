# testing/images

Images used for various CI jobs for cert-manager and related projects.
These images are all pushed to europe-west1-docker.pkg.dev/cert-manager-tests-trusted/cert-manager-infra-images/

All images are built in Prow. There is a Prow job per image in [config/jobs/testing/test-trusted.yaml](../config/jobs/testing/testing-trusted.yaml).
Most images are built using the scripts in [images/builder](./builder).

## When does a new image get built/will my change trigger a new build?

There is a Prow post-submit job that builds the image for each of the images in ./config/jobs/testing/testing-trusted.yaml.
Each of these jobs will get triggered after a change to a subdirectory in ./images, for example the job that builds new 'golang-nodejs' image will get triggered after a change to ./images/golang-nodejs, see [its '.run_if_changed' field](https://github.com/cert-manager/testing/blob/2b87fe6e34ff150042a9a776a85b3e62a20d98dc/config/jobs/testing/testing-trusted.yaml#L176).

After a PR to ./images subdirectory gets merged, you should see the Prow job that builds the new image version in https://prow.infra.cert-manager.io/.
(There is a known bug where sometimes these jobs appear as failed despite having succesfully built the image https://github.com/cert-manager/testing/issues/602)

## How do I add a new image?

Add a new postsubmit to ./config/jobs/testing/testing-trusted.yaml that watches changes to a subdirectory with your image config and runs the image build.

Add a subdirectory to ./images with the scripts to build your image. Most already existing images use the scripts in [images/builder](./builder), see for example ./golang-dind. However, you can use other means to build the image.

!! If you commit the postsubmit job that triggers new image builds as well as the new image subdirectory in a single PR, this PR will not trigger a new image build because the Prow cluster config will only be updated with the new job after the PR gets merged.
To get the image built for the first time you can either merge the postsubmit job first and then the image build scripts or trigger the first image build manually- see the section below.

## Can I trigger image build manually?

From the root of this repository run:

```bash
docker run -it -v$(pwd):/testing gcr.io/k8s-prow/mkpj --job=NAME_OF_YOUR_POSTSUBMIT_JOB--config-path=/testing/config/config.yaml --job-config-path=/testing/config/jobs/testing/testing-trusted.yaml --base-ref=master
```

This command will output a ProwJob config that you can apply to [build infra cluster](../prow/README.md)

You can then go to https://prow.infra.cert-manager.io/ to follow the build.

!! The newly triggered job will clone this repo from Github and use the image scripts from the master branch, so you cannot use it to test local changes.
### kind

[images/kind](./kind) contains a script to build our own kind image—this may be useful when needing to test against a particular version of Kubernetes for which there is no 'official' kind image yet.

To build for a new Kubernetes version, change the `KUBERNETES_VERSION` variable in `images/kind/build.sh`—­this should trigger the `post-testing-push-kind` Prow job and a `gcr.io/jetstack-build-infra-images/kind:<KUBERNETES_VERSION>` image should get built and pushed.

This image build does not use [images/builder](./builder) functionality because the image is built with `kind`, not Docker.
