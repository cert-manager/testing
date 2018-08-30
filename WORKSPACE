git_repository(
    name = "bazel_skylib",
    commit = "2169ae1c374aab4a09aa90e65efe1a3aad4e279b",
    remote = "https://github.com/bazelbuild/bazel-skylib.git",
)

load("@bazel_skylib//:lib.bzl", "versions")

versions.check(minimum_bazel_version = "0.10.0")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "1868ff68d6079e31b2f09b828b58d62e57ca8e9636edff699247c9108518570b",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.11.1/rules_go-0.11.1.tar.gz",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains(
    go_version = "1.10.2",
)

git_repository(
    name = "test_infra",
    commit = "a8cee5a60a2d9476341cf843867221a8bd18a3e8",
    remote = "https://github.com/kubernetes/test-infra.git",
)

git_repository(
    name = "io_kubernetes_build",
    commit = "84d52408a061e87d45aebf5a0867246bdf66d180",
    remote = "https://github.com/kubernetes/repo-infra.git",
)

git_repository(
    name = "io_bazel_rules_docker",
    remote = "https://github.com/bazelbuild/rules_docker.git",
    tag = "v0.5.1",
)

load("@io_bazel_rules_docker//docker:docker.bzl", "docker_repositories", "docker_pull")

docker_repositories()

git_repository(
    name = "io_bazel_rules_k8s",
    commit = "c861e4ea5a0b34e17fb682f60fa78a9c85050519",
    remote = "https://github.com/bazelbuild/rules_k8s.git",
)

new_http_archive(
    name = "yaml",
    build_file_content = """
py_library(
    name = "yaml",
    srcs = glob(["*.py"]),
    visibility = ["//visibility:public"],
)
""",
    sha256 = "592766c6303207a20efc445587778322d7f73b161bd994f227adaa341ba212ab",
    strip_prefix = "PyYAML-3.12/lib/yaml",
    urls = ["https://files.pythonhosted.org/packages/4a/85/db5a2df477072b2902b0eb892feb37d88ac635d36245a72a6a69b23b383a/PyYAML-3.12.tar.gz"],
)
