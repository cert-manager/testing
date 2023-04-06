# Image builder

The image builder is used to build test images used by ProwJobs.

In order to build an image, a simple build.yaml file is required:

```yaml
name: bazelbuild # Name of the image to be built
# Arguments that should be passed to all builds for the image
arguments:
  DOCKER_VERSION: 18.09
# Variants allow multiple images to be built in a single build step, with
# different build arguments for each build.
variants:
  "0.24.1":
    # Specify build arguments for this variant
    arguments:
        BAZEL_VERSION: 0.24.1
# Image names to be tagged and pushed
images:
- ${_REGISTRY}/${_NAME}:${_DATE_STAMP}-${_GIT_REF}-${BAZEL_VERSION}
- ${_REGISTRY}/${_NAME}:latest-${BAZEL_VERSION}
```

## Running

To build an image locally, from the root of this repository run:

```bash
$ ./images/builder/ci-runner.sh images/bazelbuild
```

### Additional options



### Built in build arguments

The builder automatically injects build variables into image builds, and makes
them available for templating in the `images` section of the `build.yaml` file.

+-------------+------------------------------------------------------+---------------------------------------+
| Name        | Description                                          | Example                               |
+-------------+------------------------------------------------------+---------------------------------------+
| _NAME       | The name of the image as specified in build.yaml     | bazelbuild                            |
| _REGISTRY   | The image registry (specified as --registry)         | eu.gcr.io/jetstack-build-infra-images |
| _DATE_STAMP | The current date stamp, useful for use in image tags | 20190407                              |
| _GIT_REF    | The current git reference of the repository          | 2ba5d19                               |
| _VARIANT    | The name of the variant being built, if any          | experimental                          |
+-------------+------------------------------------------------------+---------------------------------------+

Additionally, all global and variant-specific options will be provided to the
Docker build, and made available for templating as part of the `images` section.
