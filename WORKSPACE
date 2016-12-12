workspace(name = "wstunnel")

# Go support
git_repository(
        name = "io_bazel_rules_go",
        remote = "https://github.com/bazelbuild/rules_go.git",
        tag = "0.3.1",
    )
load("@io_bazel_rules_go//go:def.bzl", "go_repositories", "new_go_repository")

go_repositories()

new_go_repository(
    name = "org_github_go_socks5",
    commit = "e75332964ef517daa070d7c38a9466a0d687e0a5",
    importpath = "github.com/armon/go-socks5",
)

new_go_repository(
    name = "org_golang_x_net",
    commit = "b1a2d6e8c8b5fc8f601ead62536f02a8e1b6217d",
    importpath = "golang.org/x/net",
)
