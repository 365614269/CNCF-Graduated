load(
    "@envoy_build_config//:extensions_build_config.bzl",
    "LEGACY_ALWAYSLINK",
)
load(
    "//bazel:envoy_build_system.bzl",
    "envoy_cc_library",
    "envoy_package",
)

licenses(["notice"])  # Apache 2

envoy_package()

envoy_cc_library(
    name = "data_impl_lib",
    hdrs = ["data_impl.h"],
    deps = [
        "//envoy/http:filter_interface",
        "//envoy/http:header_map_interface",
    ],
)

envoy_cc_library(
    name = "inputs_lib",
    srcs = ["inputs.cc"],
    hdrs = [
        "inputs.h",
    ],
    deps = [
        "//envoy/http:filter_interface",
        "//envoy/http:header_map_interface",
        "//source/common/http:header_utility_lib",
        "//source/common/http:utility_lib",
        "@envoy_api//envoy/type/matcher/v3:pkg_cc_proto",
    ],
    alwayslink = LEGACY_ALWAYSLINK,
)

envoy_cc_library(
    name = "status_code_input_lib",
    srcs = ["status_code_input.cc"],
    hdrs = [
        "status_code_input.h",
    ],
    deps = [
        "//envoy/http:header_map_interface",
        "//source/common/http:header_utility_lib",
        "//source/common/http:utility_lib",
        "@envoy_api//envoy/type/matcher/v3:pkg_cc_proto",
    ],
)
