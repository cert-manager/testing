# Prow deployment

This directory contains the manifests used for the deployment of the Prow
cluster.

## Upgrading Prow

The Kubernetes [Prow
deployment](https://github.com/kubernetes/test-infra/tree/master/prow) is
automatically deployed, but all the other projects like Knative, Istio, and
cert-manager do the deployment manually.

Here is the process to upgrade Prow:

1. ⚠️ You must be given the role `roles/container.developer` on the
   [jetstack-build-infra](https://console.cloud.google.com/home/dashboard?project=jetstack-build-infra)
   project. You must be able to run `kubectl` commands on the
   [github-build-infra](https://console.cloud.google.com/kubernetes/clusters/details/europe-west1-b/github-build-infra/details?project=jetstack-build-infra)
   cluster.
2. Clone this repo:

   ```sh
   git clone https://github.com/jetstack/testing
   cd testing
   ```

3. Pick a build of Prow by running:

   ```sh
   % gcloud container images list-tags gcr.io/k8s-prow/deck | head
   DIGEST       TAGS                                     TIMESTAMP
   96dba717b1f3 latest,latest-root,v20210412-ed35ec0cee  2021-04-12T16:17:11
   255fe5a57fb4 v20210412-176e4b678c                     2021-04-12T15:39:17
   53107953d93e v20210412-f0c722e283                     2021-04-12T14:59:15
   f2eca760c0f9 v20210410-57fae234ba                     2021-04-10T02:55:02
   ```

   For example, let us pick the latest one. What we call the "target commit" in
   the next steps is the commit hash that appears in the image tag:

   ```sh
   v20210412-ed35ec0cee
   #         <-------->
   #        target commit
   ```

   In this example,
   [ed35ec0cee](https://github.com/kubernetes/test-infra/commit/ed35ec0cee) is
   the target commit to which you will be upgrading to (Prow does not have
   "releases").

4. Find out what is the "current commit" of the current deployment of Prow. This
   is stored in the file `prow/version`. For example:

   ```sh
   % cat prow/version
   v20200628-cc1c099dad
   #         <-------->
   #        current commit
   ```

   At this point, you know that:

   |                 | image tag            | commit         |
   | --------------- | -------------------- | -------------- |
   | current version | v20200628-cc1c099dad | [cc1c099dad][] |
   | target version  | v20210412-ed35ec0cee | [ed35ec0cee][] |

   [cc1c099dad]: https://github.com/kubernetes/test-infra/commit/cc1c099dad
   [ed35ec0cee]: https://github.com/kubernetes/test-infra/commit/ed35ec0cee

5. Open
   [ANNOUNCEMENTS.md](https://github.com/kubernetes/test-infra/blob/master/prow/ANNOUNCEMENTS.md)
   and look for anything that changed between the current commit and the target
   commit.
6. Update the file `prow/version` with your target image tag, and open a PR to
   [jetstack/infra](https://github.com/jetstack/infra). For example:

   ```diff
   diff --git a/prow/version b/prow/version
   --- a/prow/version
   +++ b/prow/version
   @@ -1 +1 @@
   -v20200628-cc1c099dad
   +v20210412-ed35ec0cee
   ```

7. Get the PR merged. Merging the PR will not do anything, we do not do rolling
   deployments.
8. Pull the latest changes from `master`. From now on, **you must be on the
   `master` branch**:

   ```sh
   git checkout master
   git pull origin master
   ```

9. Make sure you have a context in your KUBECONFIG that is called `build-infra`
   (this context name is defined in
   [print-workspace-status.sh](https://github.com/jetstack/testing/blob/master/hack/print-workspace-status.sh#L28).
   Create the `build-infra` context with:

   ```sh
   gcloud auth login
   gcloud container clusters get-credentials --project jetstack-build-infra --region europe-west1-b github-build-infra
   kubectl config rename-context gke_jetstack-build-infra_europe-west1-b_github-build-infra build-infra
   ```

10. Generate and apply the Prow manifests to the `github-build-infra` cluster:

    ```sh
    bazel run //prow/cluster:production.apply
    ```
