load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "md",
    srcs = ["md.go"],
    pure = "on",
    visibility = ["//:__subpackages__"],
    deps = [
        "@com_github_gomarkdown_markdown//:go_default_library",
        "@com_github_gomarkdown_markdown//html:go_default_library",
    ],
)
