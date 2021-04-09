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

3. Pick a commit of Prow on
   [kubernetes/test-infra](https://github.com/kubernetes/test-infra). That's the
   commit to which you will be upgrading to. We use a commit instead of a git
   tag due to Prow not having releases. For example, let's pick
   [eca83d2ac](https://github.com/kubernetes/test-infra/commit/eca83d2ac).
4. Review
   [ANNOUNCEMENTS.md](https://github.com/kubernetes/test-infra/blob/master/prow/ANNOUNCEMENTS.md)
   and look for anything that changed between the previous commit and your new
   commit.
5. Open a PR to [jetstack/infra](https://github.com/jetstack/infra) with the
   update to the `commit` field in the file [WORKSPACE](../WORKSPACE). For
   example, if you want to be upgrading from Prow
   [a8cee5a60](https://github.com/kubernetes/test-infra/commit/a8cee5a60) to
   [eca83d2ac](https://github.com/kubernetes/test-infra/commit/eca83d2ac), the
   change to WORKSPACE is:

   ```diff
   git_repository(
       name = "test_infra",
   -   commit = "a8cee5a60a2d9476341cf843867221a8bd18a3e8",
   +   commit = "eca83d2ac2b48c2732aab0c90c6eff6e564d4a21",
       remote = "https://github.com/kubernetes/test-infra.git",
   )
   ```

6. Get the PR merged. Merging the PR will not do anything, we do not do rolling
   deployments.
7. Pull the latest changes from master. From now on, you must be on the master
   branch.

8. Make sure you have a context in your KUBECONFIG that is called `build-infra`
   (this context name is defined in
   [print-workspace-status.sh](https://github.com/jetstack/testing/blob/master/hack/print-workspace-status.sh#L28).
   Create the `build-infra` context with:

   ```sh
   gcloud auth login
   gcloud container clusters get-credentials --project jetstack-build-infra --region europe-west1-b github-build-infra
   kubectl config rename-context gke_jetstack-build-infra_europe-west1-b_github-build-infra build-infra
   ```

9. Generate and apply the Prow manifests to the `github-build-infra` cluster:

   ```sh
   bazel run //prow/cluster:production.apply
   ```
