load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_embed_data")

##
## Binaries
##
go_binary(
    name = "www",
    srcs = [
        "www.go",
        "www_release.go",
    ],
    pure = "on",
    deps = [
        "//:assets",
        "@org_golang_x_crypto//acme:go_default_library",
        "@org_golang_x_crypto//acme/autocert:go_default_library",
    ],
)

go_binary(
    name = "www_debug",
    srcs = [
        "www.go",
        "www_debug.go",
    ],
    pure = "on",
    deps = ["//:assets"],
)

##
## Static assets
##
go_embed_data(
    name = "embed_assets",
    srcs = glob(["assets/**/*"]),
    package = "assets",
    var = "Asset",
)

go_library(
    name = "assets",
    srcs = [":embed_assets"],
    importpath = "github.com/BranLwyd/www/assets",
)
