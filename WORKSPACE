load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "0c10738a488239750dbf35336be13252bad7c74348f867d30c3c3e0001906096",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.23.2/rules_go-v0.23.2.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.23.2/rules_go-v0.23.2.tar.gz",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "cdb02a887a7187ea4d5a27452311a75ed8637379a1287d8eeb952138ea485f7d",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.21.1/bazel-gazelle-v0.21.1.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.21.1/bazel-gazelle-v0.21.1.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@io_bazel_rules_go//extras:embed_data_deps.bzl", "go_embed_data_dependencies")

go_rules_dependencies()

go_register_toolchains()

go_embed_data_dependencies()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

go_repository(
    name = "com_github_branlwyd_acnh_flowers",
    commit = "dc2082decf1c7fef953e34b039d69280a9e318e0",
    importpath = "github.com/BranLwyd/acnh_flowers",
)

go_repository(
    name = "com_github_gomarkdown_markdown",
    commit = "3f9352745725482bc45bab368fdd4d111ea67307",
    importpath = "github.com/gomarkdown/markdown",
)

go_repository(
    name = "org_golang_x_crypto",
    commit = "948cd5f35899cbf089c620b3caeac9b60fa08704",
    importpath = "golang.org/x/crypto",
)

go_repository(
    name = "org_golang_x_net",
    commit = "ab34263943818b32f575efc978a3d24e80b04bd7",
    importpath = "golang.org/x/net",
)

go_repository(
    name = "org_golang_x_text",
    commit = "23ae387dee1f90d29a23c0e87ee0b46038fbed0e",
    importpath = "golang.org/x/text",
)
