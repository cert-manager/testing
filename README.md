# jetstack/testing

This repository contains the configuration used for testing all jetstck projects.

It is used by [Prow](https://github.com/kubernetes/test-infra/tree/master/prow)
to provide GitHub automation to all of our repositories.

## Common tasks

* [config/](config/): Adding or modifying CI jobs (presubmits, periodics or postsubmits)
* [prow/](prow/): Updating/upgrading Prow
* [images/](images/): Creating or modifying images used during CI

## Structure

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

These images are currently manually pushed as and when required using the ad-hoc
Makefile's contained in each directory.

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
