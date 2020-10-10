load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "b725e6497741d7fc2d55fcc29a276627d10e43fa5d0bb692692890ae30d98d00",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.24.3/rules_go-v0.24.3.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.24.3/rules_go-v0.24.3.tar.gz",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "b85f48fa105c4403326e9525ad2b2cc437babaa6e15a3fc0b1dbab0ab064bc7c",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.2/bazel-gazelle-v0.22.2.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.2/bazel-gazelle-v0.22.2.tar.gz",
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
    commit = "8c8b3816f167b780c855b6412793ffd5de35ef05",
    importpath = "github.com/gomarkdown/markdown",
)

go_repository(
    name = "org_golang_x_crypto",
    commit = "7f63de1d35b0f77fa2b9faea3e7deb402a2383c8",
    importpath = "golang.org/x/crypto",
)

go_repository(
    name = "org_golang_x_net",
    commit = "dbdefad45b8998787d2fd6a2edea7fc79838207b",
    importpath = "golang.org/x/net",
)

go_repository(
    name = "org_golang_x_text",
    commit = "a8b4671254579a87fadf9f7fa577dc7368e9d009",
    importpath = "golang.org/x/text",
)
