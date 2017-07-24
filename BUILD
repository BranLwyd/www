load("@io_bazel_rules_go//go:def.bzl", "go_prefix", "go_binary", "go_library")

go_prefix("github.com/BranLwyd/www")

##
## Binaries
##
go_binary(
    name = "www",
    srcs = [
        "www.go",
        "www_release.go",
    ],
    deps = [
        "//:assets",
        "@org_golang_x_crypto//acme/autocert:go_default_library",
    ],
)

go_binary(
    name = "www_debug",
    srcs = [
        "www.go",
        "www_debug.go",
    ],
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
    cmd = "go-bindata -o $@ --nomemcopy --nocompress --pkg=assets --prefix=assets/ $(locations :assets_files)",
)

go_library(
    name = "assets",
    srcs = ["assets.go"],
)
