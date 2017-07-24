load("@io_bazel_rules_go//go:def.bzl", "go_prefix", "go_library")

go_prefix("github.com/BranLwyd/www")

filegroup(
    name = "site",
    srcs = glob(["site/**/*"]),
)

genrule(
    name = "bindata",
    srcs = [":site"],
    outs = ["bindata.go"],
    cmd = "go-bindata -o $@ -nocompress -nomemcopy -pkg=data -prefix=data/site/ $(locations :site)",
)

go_library(
    name = "go_default_library",
    srcs = [
        "bindata.go",
        "doc.go",
    ],
    visibility = ["//:__pkg__"],
)
