http_archive(
    name = "io_bazel_rules_go",
    sha256 = "f5127a8f911468cd0b2d7a141f17253db81177523e4429796e14d429f5444f5f",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.16.1/rules_go-0.16.1.tar.gz",
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "c0a5739d12c6d05b6c1ad56f2200cb0b57c5a70e03ebd2f7b87ce88cabf09c7b",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.14.0/bazel-gazelle-0.14.0.tar.gz"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

go_rules_dependencies()

go_register_toolchains()

gazelle_dependencies()

go_repository(
    name = "org_golang_x_crypto",
    commit = "74cb1d3d52f4c01cbfb44c1b50d204462f3124c7",
    importpath = "golang.org/x/crypto",
)

go_repository(
    name = "com_github_jteeuwen_go-bindata",
    commit = "a0ff2567cfb70903282db057e799fd826784d41d",
    importpath = "github.com/jteeuwen/go-bindata",
)
