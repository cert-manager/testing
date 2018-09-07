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

Prow's deployment on our build-infra cluster is completely managed via Bazel too.

Bazel will take care of building docker images for each Prow component, as well
as pushing those images to a remote repository and rolling them out to our Prow
cluster.

The code for the Prow components themselves exists in our fork of [test-infra](https://github.com/jetstack/test-infra).

After changes have been made in our *test-infra* fork, you will need to perform
a few steps in this repository to roll them out:

#### 0. Correctly configure your local KUBECONFIG

Bazel **will not** automatically configure your KUBECONFIG file to point to our
build clusters. This is by design.

In order to be able to deploy Prow itself, we must ensure our KUBECONFIG is
configured with **two** contexts with appropriate names.

* **build-infra** - this context should be configured with credentials to talk
to the cluster running **the Prow control plane**. The credentials here will be
used to deploy Prow itself.

You can configure this with:

```
$ gcloud container clusters get-credentials \
    github-build-infra \
    --zone europe-west1-b \
    --project jetstack-build-infra

$ kubectl config rename-context gke_jetstack-build-infra_europe-west1-b_github-build-infra build-infra
```

* **libvirt** - this context should be configured with credentials to talk to
the cluster running the **Prow CI jobs** (i.e. our 'libvirt' cluster).
As this cluster is currently deployed on GCE, exact details for obtaining these
credentials vary. You will need to obtain them from the Terraform project used
to deploy the cluster. This may then require some manual 'merging' of KUBECONFIG
files in order to make both contexts available in the same kubeconfig file.

The names of these two contexts is defined in `hack/print-workspace-status.sh`.
In the unlikely event you need to change them, you can do so there.

This step will likely only need to be done once, provided you do not regularly
delete your KUBECONFIG!

#### 1. Bump the test-infra version in our WORKSPACE file

First, we must bump the version of test-infra that we reference in our
[WORKSPACE](WORKSPACE) file. If you open the file, you should see a rule that
looks something like:

```python
git_repository(
    name = "test_infra",
    commit = "24b536d5e1714637e4433bacddffd9efeb1044cb",
    remote = "https://github.com/jetstack/test-infra.git",
)
```

Change the commit ref that is referenced here to the new reference in the *test-infra*
repository.

#### 2. Build, push and deploy the new Prow components

The entire build/push/deploy workflow is handled by Bazel rules defined in
[prow/](prow/).

Now that we have updated our WORKSPACE file to point to the new version of *test-infra*,
we can run the following command which will automatically roll out the new version
of Prow:

```bash
# Obtain credentials for the docker registry
$ gcloud docker -a
# Build, push and deploy the new Prow components
$ bazel run //prow/cluster:production.apply
```

After this has run to completion, the new version should be running.
This command will deploy resources to both the **build-infra** and **libvirt**
clusters as appropriate.

Please manually verify the Prow deployment is 'healthy' after rolling out changes.
We do not currently have an automated process for this, aside from Prometheus
alerting/monitoring.
You should attempt to run at least one presubmit/postsubmit/periodic to verify.

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

These images are currently manually pushed as and when required.

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
