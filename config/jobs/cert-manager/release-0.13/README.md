# Release presubmits

This directory contains presubmit definitions for previous releases of cert-manager.

We explicitly define these as the requirements for the build on master may change,
consequently breaking the build for previous release branches.

By controlling and versioning presubmits for release branches separately, we can
be sure our release branches continue to pass tests when we make changes to the
build on master.

## When to update these

When a new release branch is created and will no longer be fast-forward to HEAD
of cert-manager master, we must snapshot the
[cert-manager-presubmits.yaml](../cert-manager-presubmits.yaml) file and copying
it across to a new file in this directory.

One minor adjustment must be made to each presubmit: the branch that the presubmit
is targetting, e.g.:

```yaml
...

  - name: pull-cert-manager-unit
    always_run: true
    skip_report: false
    context: pull-cert-manager-unit
    max_concurrency: 2
    agent: kubernetes
    decorate: true
    branches:
    - master

...
```
