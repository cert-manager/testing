# Prow deployment

Currently our Prow instance is deployed manually with Bazel using the static manifests in [./cluster](./cluster).

Prow's 'control plane' is deployed to `github-build-infra` (referred to as `build-infra`) cluster in `jetstack-build-infra` project.

Prow will spin up test pods in `jetstack-build-infra-workers-gke` (also referred to as 'default') cluster in `jetstack-build-infra-gke` project and in `jetstack-build-infra-workers-trusted` (also referred to as 'trusted) cluster in `jetstack-build-infra-internal` project depending on the type of the job.

## Upgrading Prow

New images for Prow components are built upstream on all commits to [k/test-infra/prow](https://github.com/kubernetes/test-infra/tree/master/prow)

Upgrade steps:

1. Checkout the master branch of this repo. **All commands must be run from the master branch* and from the root of this repo**. You can make the version-related changes on your locally on master branch, upgrade the components in cluster using the local changes and push your changes to Git once you have verified that the upgrade worked.

1. Ensure that you have been granted `roles/container.developer` role on the
   [jetstack-build-infra](https://console.cloud.google.com/home/dashboard?project=jetstack-build-infra)
   project

2. Configure your KUBECONFIG to point at `build-infra` cluster. The context **must** be named 'build-infra'.
Bazel **will not** automatically configure your KUBECONFIG file. This is by design.

```sh
$ gcloud container clusters get-credentials \
    github-build-infra \
    --zone europe-west1-b \
    --project jetstack-build-infra

$ kubectl config rename-context gke_jetstack-build-infra_europe-west1-b_github-build-infra build-infra
```
The name of this context is defined in `hack/print-workspace-status.sh`.
In the unlikely event you need to change it, you can do so there.

3. Ensure that you can access the cluster and view Prow components, might be worth checking component logs at this point, so you are aware which warnings/errors were present already before the upgrade.

4. Find out the latest version of upstream components:

   ```sh
   % gcloud container images list-tags gcr.io/k8s-prow/deck | head
   DIGEST       TAGS                                     TIMESTAMP
   96dba717b1f3 latest,latest-root,v20210412-ed35ec0cee  2021-04-12T16:17:11
   255fe5a57fb4 v20210412-176e4b678c                     2021-04-12T15:39:17
   53107953d93e v20210412-f0c722e283                     2021-04-12T14:59:15
   f2eca760c0f9 v20210410-57fae234ba                     2021-04-10T02:55:02
   ```

5. Check the release notes.
Prow does not have semver-versioned releases, but the image tags contain the SHA of the commit from which the image was built- so you can use commit times to detemine the relevant new changes from [k/test-infra/ANNOUNCEMENTS.md](https://github.com/kubernetes/test-infra/blob/master/prow/ANNOUNCEMENTS.md)

6. Update the [./prow/version](./version) file with the selected image tag.

7. Bump the image tags in static manifests using [./prow/bump](./bump)
This tool will read the version from `./prow/version` file.

```go
go run prow/bump/main.go
```

This should have updated image tags in the static manifest files in [./prow/cluster](./cluster).

8. Apply the updated manifests to `build-infra` cluster.

```sh
bazel run //prow/cluster:production.apply
```

9. Verify the upgrade:

- Check that all `Deployment`s and `Daemonset`s are up and running and up to date

- Check Prow component pod logs for any errors

- Trigger an e2e test and see it succeed

- Ensure you can access `https://prow.build-infra.jetstack.net/` (and see logs for the tests there) and `https://triage.build-infra.jetstack.net/s/daily`

10. Commit and PR in your change


* TODO: check if that is the case and why