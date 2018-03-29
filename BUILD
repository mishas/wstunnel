package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_prefix", "go_binary")

go_prefix("github.com/mishas/websocket-tunnel")

go_binary(
    name = "client",
    srcs = ["client.go"],
    pure = "on",
    deps = [
        "@org_golang_x_net//proxy:go_default_library",
        "@org_golang_x_net//websocket:go_default_library",
    ],
)

go_binary(
    name = "server",
    srcs = ["server.go"],
    pure = "on",
    deps = [
        "@org_github_go_socks5//:go_default_library",
        "@org_golang_x_net//websocket:go_default_library",
    ],
)
