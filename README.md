# jetstack/testing

This repository contains the configuration used for testing all jetstck projects.

It is used by [Prow](https://github.com/kubernetes/test-infra/tree/master/prow)
to provide GitHub automation to all of our repositories.

## Common tasks

### Linting files in this repository

We have certain requirements on files in these repository:

* boilerplate check - we require that all files in the repository have a valid
copyright notice at the top of the file. Examples of copyright notices for
different filetypes can be seen in [hack/boilerplate](hack/boilerplate).

* bazel file check - we 'lint' our Bazel files, as well as auto-generate the
`package-srcs` and `all-srcs` targets in each.

You can run the lint checks with:

```bash
bazel test //hack/...
```

Alternatively, to test a single aspect, you can:

```bash
# Run Bazel lint checks
bazel test //hack:verify-kazel

# Run boilerplate checker
bazel test //hack:verify-boilerplate
```

### Running the Bazel linter

As noted above, we lint our Bazel files and auto-generate targets for certain
common tasks.

You can auto-lint and auto-generate these targets like so:

```bash
bazel run //hack:update-kazel
```

### Validating Prow configuration

In order to test the configuration is valid, you can run:

```
bazel test //config/...
```

This will use the test-infra 'checkconfig' tool to verify the configuration
files.

### Deploying a new version of Prow

Prow's deployment on our build-infra cluster is done manually using Bazel
scripts in ./prow/cluster.

See more detailed information about upgrading Prow in [./prow/cluster/README.md](./prow/cluster/README.md)

### Building an image and exporting to your local Docker daemon

Each directory under `images/` and `legacy/images` contains a Bazel build file
defining how each image should be built.

You can build these images and store them within your local docker daemon by
running:

```bash
$ bazel run //images/bazelbuild
INFO: Analysed target //images/bazelbuild:bazelbuild (1 packages loaded).
INFO: Found 1 target...
Target //images/bazelbuild:bazelbuild up-to-date (nothing to build)
INFO: Elapsed time: 0.783s, Critical Path: 0.08s
INFO: 0 processes.
INFO: Build completed successfully, 1 total action
INFO: Build completed successfully, 1 total action
Loaded image ID: sha256:3c6a6d4f8f7c760670825a52475029dbc0da333eebed5472ece60fdd6ed51949
Tagging 3c6a6d4f8f7c760670825a52475029dbc0da333eebed5472ece60fdd6ed51949 as eu.gcr.io/jetstack-build-infra-images/bazelbuild:v20180907-8793fc5-0.16.1
```

This may take a few minutes depending on the state of your Bazel & Docker cache.

### Pushing a docker image to the image repository

Bazel is used to *push* built docker images to the remote registry.

Each images directory exposes a `push` rule in its BUILD.bazel file that can be
used to push images automatically.

This push target **will not** handle authentication with the remote registry for
you. You should ensure your Docker client is already authenticated using gcloud.

For example, to build and push the `images/bazelbuild` image:

```bash
# Obtain credentials for the docker registry
$ gcloud docker -a
# Build (if required) and push the docker image
$ bazel run //images/bazelbuild:push
...
```

Again, this may take a few minutes depending on the state of your Bazel and
Docker cache.

The docker repository that will be pushed to is defined in `hack/print-workspace-status.sh`.
If you want to push to a custom repository, you will need to edit this file
manually.
In future, we may allow this to be overridden using environment variables or
build arguments passed to `bazel run`.

---

## Structure

* [config/](config/): Adding or modifying CI jobs (presubmits, periodics or postsubmits)
* [prow/](prow/): Updating/upgrading Prow
* [images/](images/): Creating or modifying images used during CI

### hack/

This contains a bazel build file and support scripts used to verify aspects
of the repository.

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


### legacy/

Prow supports two modes for configuring jobs - 'decorated' and 'bootstrap'.

The decorated mode uses init containers and a sidecar to perform job 'utility'
functions, such as uploading logs to GCS and cloning the repo you are testing
at the correct revision. This is a newer approach, with a few limitations in
the amount of build metadata can be displayed. It is recommended to be used
going forward.

The bootstrap approach relies on a Python script in this repository, under [legacy/bootstrap](legacy/bootstrap).

A number of our jobs still rely on this 'bootstrap' approach, and as such we
maintain a copy of all required files within this configuration repository.

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

1. Remove Bazel presets from job config, so it doesn't look for Bazel cache creds
2. Run `./prow/pj-on-kind.sh pull-cert-manager-upgrade-v1-21`
3. Pass some cert-manager PR number when requested. This will be checked out.
4. Pass 'empty' for any storage volumes when requested.
5. Retrieve kubeconfig for the kind cluster `kind get kubeconfig --name mkpod` and set KUBECONFIG
6. `kubectl get pods` - to get the name of the pod that is running the test
7. `kubectl logs <pod-name> -c test -f` stream the logs