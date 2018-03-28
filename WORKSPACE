workspace(name = "wstunnel")

# Go support
git_repository(
    name = "io_bazel_rules_go",
    remote = "https://github.com/bazelbuild/rules_go.git",
    tag = "0.10.1",
)
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")
go_rules_dependencies()
go_register_toolchains()

go_repository(
    name = "org_github_go_socks5",
    commit = "e75332964ef517daa070d7c38a9466a0d687e0a5",
    importpath = "github.com/armon/go-socks5",
)
