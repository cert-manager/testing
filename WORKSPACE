workspace(name = "io_jetstack_testing")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_file")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "new_git_repository")

git_repository(
    name = "bazel_skylib",
    remote = "https://github.com/bazelbuild/bazel-skylib.git",
    tag = "0.6.0",
)

load("@bazel_skylib//:lib.bzl", "versions")

versions.check(minimum_bazel_version = "0.15.0")

http_archive(
    name = "io_bazel_rules_go",
    urls = ["https://github.com/bazelbuild/rules_go/releases/download/0.16.5/rules_go-0.16.5.tar.gz"],
    sha256 = "7be7dc01f1e0afdba6c8eb2b43d2fa01c743be1b9273ab1eaf6c233df078d705",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()
go_register_toolchains(
    go_version = "1.11.4",
)

git_repository(
    name = "test_infra",
    commit = "d8438622095051e201f8991099a3ba513e579661",
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
    tag = "v0.6.0",
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


http_archive(
    name = "build_bazel_rules_typescript",
    strip_prefix = "rules_typescript-0.22.0",
    url = "https://github.com/bazelbuild/rules_typescript/archive/0.22.0.zip",
)

# Fetch our Bazel dependencies that aren't distributed on npm
load("@build_bazel_rules_typescript//:package.bzl", "rules_typescript_dependencies")

rules_typescript_dependencies()

# Setup TypeScript toolchain
load("@build_bazel_rules_typescript//:defs.bzl", "ts_setup_workspace")

git_repository(
    name = "build_bazel_rules_nodejs",
    remote = "https://github.com/bazelbuild/rules_nodejs.git",
    tag = "0.16.4",
)

load("@build_bazel_rules_nodejs//:defs.bzl", "node_repositories", "yarn_install")

node_repositories(package_json = ["@test_infra//:package.json"])

yarn_install(
    name = "npm",
    package_json = "@test_infra//:package.json",
    yarn_lock = "@test_infra//:yarn.lock",
)

load("//def:container.bzl", "container_dockerfile")

container_dockerfile(
    name = "bazelbuild",
    build_matrix = {
        "BAZEL_VERSION": [
            "0.16.1",
            "0.18.0",
            "0.21.0",
        ],
    },
    dockerfile = "//images/bazelbuild:Dockerfile",
)

container_dockerfile(
    name = "alpine",
    dockerfile = "//images/alpine:Dockerfile",
)

container_dockerfile(
    name = "katacoda-lint",
    dockerfile = "//images/katacoda-lint:Dockerfile",
)

container_dockerfile(
    name = "gcloud-in-go",
    build_args = {
        "GCLOUD_VERSION": "163.0.0",
    },
    dockerfile = "//legacy/images/gcloud-in-go:Dockerfile",
)

container_dockerfile(
    name = "minikube-in-go",
    build_args = {
        "BAZEL_VERSION": "0.18.0",
        "GCLOUD_VERSION": "163.0.0",
    },
    build_matrix = {
        "KUBERNETES_VERSION": ["v1.9.6", "v1.8.10", "v1.7.15"],
    },
    dockerfile = "//legacy/images/minikube-in-go:Dockerfile",
)

container_dockerfile(
    name = "tarmak-ruby",
    build_args = {
        "GCLOUD_VERSION": "206.0.0",
        "GCLOUD_HASH": "d39293914b2e969bfe18dd19eb77ba96d283995f8cf1e5d7ba6ac712a3c9479a",
    },
    build_matrix = {
        "RUBY_VERSION": ["2.4.4"],
    },
    dockerfile = "//legacy/images/tarmak/ruby:Dockerfile",
)

container_dockerfile(
    name = "tarmak-sphinx-docs",
    build_args = {
        "GCLOUD_VERSION": "178.0.0",
        "GCLOUD_HASH": "2e0bbbf81c11164bf892cf0b891751ba4e5172661eff907ad1f7fc0b6907e296",
    },
    dockerfile = "//legacy/images/tarmak/sphinx-docs:Dockerfile",
)

# Dependencies used for the cert-manager e2e image
http_file(
   name = "golang",
   urls = ["https://dl.google.com/go/go1.11.linux-amd64.tar.gz"],
   sha256 = "b3fcf280ff86558e0559e185b601c9eade0fd24c900b4c63cd14d1d38613e499",
)

http_file(
    name = "dep_linux",
    executable = 1,
    urls = ["https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64"],
    sha256 = "287b08291e14f1fae8ba44374b26a2b12eb941af3497ed0ca649253e21ba2f83"
)

new_http_archive(
    name = "helm_linux",
    sha256 = "0fa2ed4983b1e4a3f90f776d08b88b0c73fd83f305b5b634175cb15e61342ffe",
    urls = ["https://storage.googleapis.com/kubernetes-helm/helm-v2.10.0-linux-amd64.tar.gz"],
    build_file_content = """
filegroup(
    name = "helm",
    srcs = [
        "linux-amd64/helm",
    ],
    visibility = ["//visibility:public"],
)
""",
)
