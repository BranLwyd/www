load("@io_bazel_rules_go//go:def.bzl", "go_prefix", "go_binary")

go_prefix("github.com/BranLwyd/www")

go_binary(
    name = "www",
    srcs = [
        "www.go",
        "www_release.go",
    ],
    deps = [
        "//data:go_default_library",
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
        "//data:go_default_library",
    ],
)
