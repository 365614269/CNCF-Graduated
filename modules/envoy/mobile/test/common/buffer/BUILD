load("@envoy//bazel:envoy_build_system.bzl", "envoy_cc_test", "envoy_mobile_package")

licenses(["notice"])  # Apache 2

envoy_mobile_package()

envoy_cc_test(
    name = "bridge_fragment_test",
    srcs = ["bridge_fragment_test.cc"],
    repository = "@envoy",
    deps = [
        "//library/common/buffer:bridge_fragment_lib",
        "//library/common/types:c_types_lib",
        "@envoy//source/common/buffer:buffer_lib",
    ],
)
