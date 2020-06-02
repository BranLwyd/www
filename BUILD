load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_embed_data", "go_library")
load("//tools:md.bzl", "convert_md_to_html")

##
## Binaries
##
go_binary(
    name = "www",
    srcs = [
        "www.go",
        "www_flower.go",
        "www_release.go",
    ],
    pure = "on",
    deps = [
        "//:assets",
        "@com_github_branlwyd_acnh_flowers//:flower",
        "@org_golang_x_crypto//acme:go_default_library",
        "@org_golang_x_crypto//acme/autocert:go_default_library",
    ],
)

go_binary(
    name = "www_debug",
    srcs = [
        "www.go",
        "www_debug.go",
        "www_flower.go",
    ],
    pure = "on",
    deps = [
        "//:assets",
        "@com_github_branlwyd_acnh_flowers//:flower",
    ],
)

##
## Static assets
##
convert_md_to_html(
    name = "html_assets",
    srcs = glob(["assets/pages/*"]),
    template = "assets/template.html",
)

go_embed_data(
    name = "embed_html_assets",
    srcs = [":html_assets"],
    flatten = True,
    package = "assets",
    var = "Page",
)

go_embed_data(
    name = "embed_static_assets",
    srcs = glob(
        ["assets/**/*"],
        exclude = ["assets/pages/*"],
    ),
    package = "assets",
    var = "Static",
)

go_library(
    name = "assets",
    srcs = [
        ":embed_html_assets",
        ":embed_static_assets",
    ],
    importpath = "github.com/BranLwyd/www/assets",
)
