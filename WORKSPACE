git_repository(
    name = "io_bazel_rules_go",
    remote = "https://github.com/bazelbuild/rules_go.git",
    tag = "0.5.5",
)

load("@io_bazel_rules_go//go:def.bzl", "go_repositories", "go_repository")

go_repository(
    name = "org_golang_x_crypto",
    commit = "51714a8c4ac1764f07ab4127d7f739351ced4759",
    importpath = "golang.org/x/crypto",
)

go_repository(
    name = "com_github_jteeuwen_go-bindata",
    commit = "a0ff2567cfb70903282db057e799fd826784d41d",
    importpath = "github.com/jteeuwen/go-bindata",
)

go_repositories()
