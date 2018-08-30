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
    commit = "a62c9b4a9dc6765256de49a9db12845491f54a1d",
    remote = "https://github.com/jetstack/test-infra.git",
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

docker_pull(
    name = "alpine-base",
    # 0.1 as of 2017/11/29
    digest = "sha256:317d39ece9dd09992fa81236964be3f3919b940f42e3143379dd66e4af930f3a",
    registry = "gcr.io",
    repository = "k8s-prow/alpine",
)

docker_pull(
    name = "git-base",
    # 0.2 as of 2018/05/10
    digest = "sha256:3eaeff9a2c35a50c3a0af7ef7cf26ea73e6fd966f54ef3dfe79d4ffb45805112",
    registry = "gcr.io",
    repository = "k8s-prow/git",
)

docker_pull(
    name = "python",
    digest = "sha256:8bfeec8f8ba3aaeea918a0198f4b1c7c9b2b39e26f399a7173229dfcef76fc1f",
    registry = "index.docker.io",
    repository = "library/python",
    tag = "2.7.14-jessie",
)

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_repositories = "repositories",
)

_go_repositories()

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
