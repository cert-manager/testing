workspace(name = "io_jetstack_testing")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_file")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "new_git_repository")

git_repository(
    name = "bazel_skylib",
    remote = "https://github.com/bazelbuild/bazel-skylib.git",
    commit = "f83cb8dd6f5658bc574ccd873e25197055265d1c",
    shallow_since = "1543273402 -0500",
)

load("@bazel_skylib//lib:versions.bzl", "versions")

versions.check(minimum_bazel_version = "0.23.0")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.22.1/rules_go-v0.22.1.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.22.1/rules_go-v0.22.1.tar.gz",
    ],
    sha256 = "e6a6c016b0663e06fa5fccf1cd8152eab8aa8180c583ec20c872f4f9953a7ac5",
)

http_archive(
    name = "bazel_gazelle",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/v0.20.0/bazel-gazelle-v0.20.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.20.0/bazel-gazelle-v0.20.0.tar.gz",
    ],
    sha256 = "d8c45ee70ec39a57e7a05e5027c32b1576cc7f16d9dd37135b0eddde45cf1b10",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

http_archive(
    name = "io_k8s_repo_infra",
    strip_prefix = "repo-infra-0.0.2",
    sha256 = "774e160ba1a2a66a736fdc39636dca799a09df015ac5e770a46ec43487ec5708",
    urls = [
        "https://github.com/kubernetes/repo-infra/archive/v0.0.2.tar.gz",
    ],
)

load("@io_k8s_repo_infra//:load.bzl", "repositories")

repositories()

load("@io_k8s_repo_infra//:repos.bzl", "configure", repo_infra_go_repositories = "go_repositories")

configure() # see repos.bzl for optional args
repo_infra_go_repositories()

git_repository(
    name = "io_bazel_rules_docker",
    remote = "https://github.com/bazelbuild/rules_docker.git",
    commit = "3772262910d1ac63563e5f1758f07df1f7857442",
    shallow_since = "1579194217 -0800",
)

load("@io_bazel_rules_docker//container:container.bzl", "container_pull")

# Used by label_sync
container_pull(
    name = "distroless-base",
    # latest circa 2017/11/29
    digest = "sha256:bef8d030c7f36dfb73a8c76137616faeea73ac5a8495d535f27c911d0db77af3",
    registry = "gcr.io",
    repository = "distroless/base",
)

container_pull(
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

# This requires rules_docker to be fully instantiated before
# it is pulled in.
# Download the rules_k8s repository at release v0.3.1
http_archive(
    name = "io_bazel_rules_k8s",
    sha256 = "cc75cf0d86312e1327d226e980efd3599704e01099b58b3c2fc4efe5e321fcd9",
    strip_prefix = "rules_k8s-0.3.1",
    urls = ["https://github.com/bazelbuild/rules_k8s/releases/download/v0.3.1/rules_k8s-v0.3.1.tar.gz"],
)

load("@io_bazel_rules_k8s//k8s:k8s.bzl", "k8s_repositories")

k8s_repositories()

load("@io_bazel_rules_k8s//k8s:k8s_go_deps.bzl", k8s_go_deps = "deps")

k8s_go_deps()

http_archive(
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

go_repository(
    name = "in_gopkg_yaml_v2",
    commit = "5420a8b6744d3b0345ab293f6fcba19c978f1183",
    remote = "https://github.com/go-yaml/yaml",
    vcs = "git",
    importpath = "gopkg.in/yaml.v2",
)

git_repository(
    name = "test_infra",
    commit = "4d31f63924b8eb14211f19a2722125b8fa0040c9",
    remote = "https://github.com/kubernetes/test-infra.git",
)
