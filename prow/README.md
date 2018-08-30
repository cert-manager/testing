# Prow deployment

This directory contains manifests used for the deployment of the Prow cluster.

## Updating Prow

The core Prow components are automatically built from our own fork of the test-infra
repository, and applied to our production build cluster.

In order to upgrade Prow to a new version, you will first need to change the
[WORKSPACE](../WORKSPACE) file in the root of this repository to reference
the desired revision.

For example, you should change the `commit` here appropriately:

```
git_repository(
    name = "test_infra",
    commit = "a8cee5a60a2d9476341cf843867221a8bd18a3e8",
    remote = "https://github.com/kubernetes/test-infra.git",
)
```

Once this is done, you can use Bazel to build, push and deploy the relevant new
images:

```
$ bazel run //prow/cluster:production.apply
```

We do *not* currently automate the roll-out of the manifests in this repository.
This means that someone with privileged access must run the `production.apply`
job.

In order to find connection details for the 'build-infra' and 'libvirt' clusters,
you will need to ensure you have two contexts already correctly configured.
Namely, `build-infra` and `libvirt`. Bazel will use the contexts with these names
to apply changes to Prow.

You can see where these context names are hardcoded in the [hack/print-workspace-status.sh](hack/print-workspace-status.sh)
file.
