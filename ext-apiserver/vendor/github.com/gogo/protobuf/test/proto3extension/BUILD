load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["proto3ext.pb.go"],
    importpath = "github.com/gogo/protobuf/test/proto3extension",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/gogo/protobuf/proto:go_default_library",
        "//vendor/github.com/gogo/protobuf/protoc-gen-gogo/descriptor:go_default_library",
    ],
)

proto_library(
    name = "proto3extension_proto",
    srcs = ["proto3ext.proto"],
    visibility = ["//visibility:public"],
    deps = ["@com_google_protobuf//:descriptor_proto"],
)

go_proto_library(
    name = "proto3extension_go_proto",
    importpath = "github.com/gogo/protobuf/test/proto3extension",
    proto = ":proto3extension_proto",
    visibility = ["//visibility:public"],
    deps = ["@com_github_golang_protobuf//ptypes/descriptor:go_default_library"],
)
