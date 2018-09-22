load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

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
    deps = [
        "//:assets",
    ],
)

##
## Static assets
##
filegroup(
    name = "assets_files",
    srcs = glob(["assets/**/*"]),
)

genrule(
    name = "assets_go",
    srcs = [":assets_files"],
    outs = ["assets.go"],
    cmd = "$(location @com_github_jteeuwen_go-bindata//go-bindata) -o $@ --nomemcopy --nocompress --pkg=assets --prefix=assets/ $(locations :assets_files)",
    tools = ["@com_github_jteeuwen_go-bindata//go-bindata"],
)

go_library(
    name = "assets",
    srcs = ["assets.go"],
    importpath = "github.com/BranLwyd/www/assets",
)
