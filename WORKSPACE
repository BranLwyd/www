load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "142dd33e38b563605f0d20e89d9ef9eda0fc3cb539a14be1bdb1350de2eda659",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.22.2/rules_go-v0.22.2.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.22.2/rules_go-v0.22.2.tar.gz",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "d8c45ee70ec39a57e7a05e5027c32b1576cc7f16d9dd37135b0eddde45cf1b10",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/v0.20.0/bazel-gazelle-v0.20.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.20.0/bazel-gazelle-v0.20.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

go_repository(
    name = "org_golang_x_crypto",
    commit = "0ec3e9974c59449edd84298612e9f16fa13368e8",
    importpath = "golang.org/x/crypto",
)

go_repository(
    name = "org_golang_x_net",
    commit = "d3edc9973b7eb1fb302b0ff2c62357091cea9a30",
    importpath = "golang.org/x/net",
)

go_repository(
    name = "org_golang_x_text",
    commit = "06d492aade888ab8698aad35476286b7b555c961",
    importpath = "golang.org/x/text",
)
