load(
    "@envoy//bazel:envoy_build_system.bzl",
    "envoy_cc_library",
    "envoy_mobile_package",
    "envoy_proto_library",
)

licenses(["notice"])  # Apache 2

envoy_mobile_package()

envoy_proto_library(
    name = "filter",
    srcs = ["filter.proto"],
    deps = [
        "@envoy_api//envoy/config/common/matcher/v3:pkg",
    ],
)

envoy_cc_library(
    name = "test_read_filter_lib",
    srcs = ["filter.cc"],
    hdrs = ["filter.h"],
    repository = "@envoy",
    deps = [
        "filter_cc_proto",
        "@envoy//source/common/http:utility_lib",
        "@envoy//source/common/stream_info:stream_info_lib",
        "@envoy//source/extensions/filters/http/common:pass_through_filter_lib",
    ],
)

envoy_cc_library(
    name = "config",
    srcs = ["config.cc"],
    hdrs = ["config.h"],
    repository = "@envoy",
    deps = [
        ":test_read_filter_lib",
        "@envoy//source/extensions/filters/http/common:factory_base_lib",
    ],
)
