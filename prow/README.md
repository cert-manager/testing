# Prow deployment

Currently our Prow instance is deployed manually with Bazel using the static manifests in [./cluster](./cluster).

Prow's 'control plane' is deployed to `prow-trusted` (referred to as `trusted`) cluster in `cert-manager-tests-trusted` project (defined in [cert-manager/infrastructure](https://github.com/cert-manager/infrastructure/blob/7b45ed95c68919c1d3cb14b8ff35fa5de46275be/gcp/clusters.tf#L3-L24)).

Prow will spin up test pods in `prow-untrusted` (also referred to as 'default') cluster in `cert-manager-tests-untrusted` project and in `prow-trusted` (also referred to as `trusted`) cluster in `cert-manager-tests-trusted` project depending on the type of the job.

The separation between 'trusted' and 'default' cluster allows us to use `ProwJob`s to perform actions that require authentication to other parts of our infrastructure (i.e push images to GCR) and at the same time protects us from a possible attack where after a maintainer has labelled a PR with 'ok-to-test', a change is made to the PR code that attacks some part of the infrastructure, i.e attempts to read `Secret`s in the cluster.
This protection works because all jobs that run in the 'trusted' cluster are periodics or postsubmit jobs- so they would not run in between a PR being 'ok-to-test'-ed and approved and merged. It is therefore important that we do not add presubmit jobs to the 'trusted' cluster.

## Upgrading Prow

New images for Prow components are built upstream on all commits to [k/kubernetes-sigs/prow](https://github.com/kubernetes-sigs/prow/tree/main)

Upgrade steps:

1. Checkout the master branch of this repo. **All commands must be run from the master branch* and from the root of this repo**. You can make the version-related changes on your locally on master branch, upgrade the components in cluster using the local changes and push your changes to Git once you have verified that the upgrade worked.

1. Ensure that you have been granted `roles/container.developer` role on the
   [cert-manager-tests-trusted](https://console.cloud.google.com/home/dashboard?project=cert-manager-tests-trusted)
   project (see [cert-manager/infrastructure](https://github.com/cert-manager/infrastructure/blob/7b45ed95c68919c1d3cb14b8ff35fa5de46275be/gcp/variables.tf#L3-L9))

2. Configure your KUBECONFIG to point at the `prow-trusted` cluster:

```sh
$ gcloud container clusters get-credentials \
   prow-trusted \
   --zone europe-west1-b \
   --project cert-manager-tests-trusted
```

3. Ensure that you can access the cluster and view Prow components, might be worth checking component logs at this point, so you are aware which warnings/errors were present already before the upgrade.

4. Checkout the autobump PR in the testing repo or manually update/ edit the versions of used images or any of the YAML in `./prow/cluster`.

5. Check the release notes.
Prow does not have semver-versioned releases, but the image tags contain the SHA of the commit from which the image was built- so you can use commit times to detemine the relevant new changes from [k/kubernetes-sigs/prow/site/content/en/docs/announcements.md](https://github.com/kubernetes-sigs/prow/blob/main/site/content/en/docs/announcements.md)

7. Review the difference between the local manifests and the live resources in the `build-infra` cluster.

```sh
cd prow
make diff-prow
```

9. Apply the updated manifests to `build-infra` cluster.

```sh
cd prow
make deploy-prow
```

10. Verify the upgrade:

- Check that all `Deployment`s and `Daemonset`s are up and running and up to date

- Check Prow component pod logs for any errors

- Trigger an e2e test and see it succeed

- Ensure you can access `https://prow.infra.cert-manager.io/` (and see logs for the tests there) and `https://triage.infra.cert-manager.io/s/daily`

11. Commit and PR in your change
