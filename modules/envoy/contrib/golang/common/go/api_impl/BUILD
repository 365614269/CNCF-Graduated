load("@io_bazel_rules_go//go:def.bzl", "go_library")

licenses(["notice"])  # Apache 2

go_library(
    name = "api_impl",
    srcs = [
        "api.h",
        "capi_impl.go",
    ],
    cgo = True,
    clinkopts = select({
        "@io_bazel_rules_go//go/platform:android": [
            "-Wl,-unresolved-symbols=ignore-all",
        ],
        "@io_bazel_rules_go//go/platform:darwin": [
            "-Wl,-undefined,dynamic_lookup",
        ],
        "@io_bazel_rules_go//go/platform:ios": [
            "-Wl,-undefined,dynamic_lookup",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "-Wl,-unresolved-symbols=ignore-all",
        ],
        "//conditions:default": [],
    }),
    importpath = "github.com/envoyproxy/envoy/contrib/golang/common/go/api_impl",
    visibility = ["//visibility:public"],
    deps = [
        "//contrib/golang/common/go/api",
    ],
)
