load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "trap_others.go",
    ] + select({
        "@io_bazel_rules_go//go/platform:windows_amd64": [
            "trap_windows.go",
            "trap_windows_1.4.go",
        ],
        "//conditions:default": [],
    }),
    importpath = "github.com/inconshreveable/mousetrap",
    visibility = ["//visibility:public"],
)
