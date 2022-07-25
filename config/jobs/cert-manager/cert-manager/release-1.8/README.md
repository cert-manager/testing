# release-1.8 Prow Specs

release-1.8 is slightly unique in that it also has a selection of hand-rolled tests, as well as automatically generated tests.

This is because 1.8 was the last version to be released with Bazel still being a part of the process. We still need to use bazel
for some tests and to ensure bazel works for the build.

Rather than write generators for bazel tests, we hand roll those and maintain them separately.
