# Legacy bootstrap-based job configuration

> The files in this directory are now abandoned and should not be used.
> They remain here for a shortwhile in order to support existing jobs that
> have not yet migrated.

This directory contains the supporting files needed for the legacy 'bootstrap'
ProwJob setup.

It provides a slim-down version of what is contained in the upstream test-infra
repository, containing only the files required for running jobs using bootstrap.py.

## Usage

You can run jobs using the bootstrap.py script like so:

```
$ bootstrap/bootstrap.py \
    --job=pull-cert-manager-verify \
    --repo=github.com/cert-manager/cert-manager=master \
    --scenario=execute \
    -- \
    make verify
```

CI scripts/images should be updated to clone this repo instead of cloning the
entirety of `test-infra`.

## Structure

### bootstrap/

This directory contains the actual bootstrap script, taken from the `jenkins/`
directory in the upstream test-infra

### jobs/

This directory contains the old-style `config.json` file for defining jobs and
their scenario mappings.

### scenarios/

This directory contains numerous python scripts that can be used to bootstrap
test environments. Similar to the other directories, it has also been taken from
the `test-infra` repository and serves a similar purpose.
