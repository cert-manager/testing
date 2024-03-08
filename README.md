# cert-manager/testing

This repository contains the configuration used for testing all jetstck projects.

It is used by [Prow](https://github.com/kubernetes/test-infra/tree/master/prow)
to provide GitHub automation to all of our repositories.

## Common tasks

### Linting files in this repository

We have certain requirements on files in these repository:

* boilerplate check - we require that all files in the repository have a valid
copyright notice at the top of the file. Examples of copyright notices for
different filetypes can be seen in [hack/boilerplate](hack/boilerplate).

You can run the lint checks with:

```bash
make verify
```

### Validating Prow configuration

In order to test the configuration is valid, you can run:

```
make local-checkconfig
```

This will use the test-infra 'checkconfig' tool to verify the configuration
files.

### Deploying a new version of Prow

Prow's deployment on our build-infra cluster is done manually using a Makefile in ./prow/cluster.

See more detailed information about upgrading Prow in [./prow/cluster/README.md](./prow/cluster/README.md)

### Building an image and exporting to your local Docker daemon

Each directory under `images/` contains a configuration file that
define how each image should be built.

You can build these images and store them within your local docker daemon by
running:

```bash
$ ./images/image-builder-script/builder.sh images/golang-dind
./images/image-builder-script/builder.sh images/golang-aws
WARNING: GOOGLE_APPLICATION_CREDENTIALS not set
Executing builder...
2023/04/07 16:31:51 --confirm is set to false, not pushing images
...
```

This may take a few minutes depending on the state of your Docker cache.

### Pushing a docker image to the image repository

⚠️ WARNING: You're unlikely to have permissions to be able to push images to GCR locally. If you're simply
looking to update an image, a [workload](https://github.com/cert-manager/testing/blob/365d570125e751a7d9aac4148d8c0ef23e42168c/config/jobs/testing/testing-postsubmits-trusted.yaml#L76)
in prow will build and push the image for you when your PR with the changes is merged.

builder.sh can also be used to *push* built docker images to the remote registry.

This push target **will not** handle authentication with the remote registry for
you. You should ensure your Docker client is already authenticated using gcloud.

For example, to build and push the `images/golang-aws` image:

```bash
# Obtain credentials for the docker registry
$ gcloud docker -a
# Build (if required) and push the docker image
$ ./images/image-builder-script/builder.sh images/golang-aws --confirm=true
...
```

If you want to push to a custom repository, you can use the `--registry` flag.

---

## Structure

* [config/](config/): Adding or modifying CI jobs (presubmits, periodics or postsubmits)
* [prow/](prow/): Updating/upgrading Prow
* [images/](images/): Creating or modifying images used during CI

### hack/

This contains support scripts used to verify aspects of the repository.

### config/

The config directory contains the actual Prow configuration files: `config.yaml`
and `plugins.yaml`, as well as job configuration for presubmits, periodics and
postsubmits.

Pull requests can be submitted to this directory in order to modify how tests
are run.

Once your changes have been merged, Prow itself will automatically update its
configuration to reflect what is in the repository.

### images/

This directory contains image defintions for images used as part of Prow jobs.

New images will be built and pushed on changes to the relevant files (i.e
Dockerfile for the image).

### Debugging e2e tests run with Prow

-  For each e2e test run, Prow will create a new `ProwJob` custom resource in
   `build-infra` cluster. For the actual test, a pod will be spun up in
   `build-infra-workers` cluster in `test-pods` namespace. You can find the pod's
   name from the `ProwJob`'s yaml `kubectl get prowjob <prowjob-name> -ojsonpath='{.status.pod_name}'`

- When debugging a periodic Prow test, a new test run can be triggered by
  deleting the latest `ProwJob` for that test

- The image used for the test container has bash, so a running test can be
  easily debugged by execing the container `kubectl exec -it  <pod-name> -ctest
  -ntest-pods -- bash`

- When execed to test container, you can find tools such as `kubectl`, `kind`, `helm`,
  `jq` in `~/bazel-out/k8-fastbuild/bin/hack/bin/`. The current kube context will
  already be that of the kind cluster that runs the e2e tests

## Creating new Prowjobs

See documentation for ProwJobs in [k/test-infra](https://github.com/kubernetes/test-infra/blob/master/prow/jobs.md).

### Testing locally

ProwJobs can be tested locally by running the (interactive) `./prow/pj-on-kind.sh` script.
This script will spin up a local KIND cluster and create a new ProwJob instance for which there will be a Pod created that will be running the actual test.

See [documentation in k/test-infra](https://github.com/kubernetes/test-infra/blob/master/prow/build_test_update.md#How-to-test-a-ProwJob) for how the script works.

An example of running `pull-cert-manager-upgrade-v1-21` job locally:

1. Run `./prow/pj-on-kind.sh pull-cert-manager-upgrade-v1-21`
2. Pass some cert-manager PR number when requested. This will be checked out.
3. Pass 'empty' for any storage volumes when requested.
4. Retrieve kubeconfig for the kind cluster `kind get kubeconfig --name mkpod` and set KUBECONFIG
5. `kubectl get pods` - to get the name of the pod that is running the test
6. `kubectl logs <pod-name> -c test -f` stream the logs
