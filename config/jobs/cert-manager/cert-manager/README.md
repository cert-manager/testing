# cert-manager Prow Specs

cert-manager prow jobs are defined based on the branch they're tested against, and only supported versions are tested.

That means that after a new major version of cert-manager is released, tests for now-deprecated versions should be manually
removed.

## Generating Tests

Tests are generated using [`cmrel`](https://github.com/cert-manager/release).

```console
go install github.com/cert-manager/release/cmd/cmrel@latest
cmrel generate-prow --help
```

For example, to generate the latest tests for master:

```console
cmrel generate-prow --branch=master
```
