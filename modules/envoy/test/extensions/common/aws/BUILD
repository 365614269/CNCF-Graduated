load(
    "//bazel:envoy_build_system.bzl",
    "envoy_cc_mock",
    "envoy_cc_test",
    "envoy_package",
    "envoy_select_boringssl",
)

licenses(["notice"])  # Apache 2

envoy_package()

envoy_cc_mock(
    name = "aws_mocks",
    srcs = ["mocks.cc"],
    hdrs = ["mocks.h"],
    rbe_pool = "6gig",
    deps = [
        "//envoy/http:message_interface",
        "//source/common/http:message_lib",
        "//source/extensions/common/aws:aws_cluster_manager_lib",
        "//source/extensions/common/aws:credential_provider_chains_lib",
        "//source/extensions/common/aws:credentials_provider_base_lib",
        "//source/extensions/common/aws:credentials_provider_interface",
        "//source/extensions/common/aws:metadata_fetcher_lib",
        "//source/extensions/common/aws:signer_interface",
        "//source/extensions/common/aws/credential_providers:iam_roles_anywhere_x509_credentials_provider_lib",
        "//source/extensions/common/aws/signers:iam_roles_anywhere_sigv4_signer_impl_lib",
        "//source/extensions/common/aws/signers:sigv4a_signer_impl_lib",
    ],
)

envoy_cc_test(
    name = "credential_provider_chains_test",
    srcs = ["credential_provider_chains_test.cc"],
    rbe_pool = "6gig",
    deps = [
        "//source/extensions/common/aws:credential_provider_chains_lib",
        "//source/extensions/common/aws/credential_providers:inline_credentials_provider_lib",
        "//test/extensions/common/aws:aws_mocks",
        "//test/mocks/api:api_mocks",
        "//test/mocks/server:factory_context_mocks",
        "//test/mocks/upstream:cluster_manager_mocks",
        "//test/test_common:environment_lib",
        "//test/test_common:test_runtime_lib",
    ],
)

envoy_cc_test(
    name = "aws_cluster_manager_test",
    srcs = ["aws_cluster_manager_test.cc"],
    rbe_pool = "6gig",
    deps = [
        "//source/extensions/common/aws:aws_cluster_manager_lib",
        "//test/mocks/server:server_factory_context_mocks",
        "@envoy_api//envoy/config/cluster/v3:pkg_cc_proto",
    ],
)

envoy_cc_test(
    name = "utility_test",
    srcs = ["utility_test.cc"],
    copts = envoy_select_boringssl(
        [
            "-DENVOY_SSL_FIPS",
        ],
    ),
    rbe_pool = "6gig",
    deps = [
        "//source/extensions/common/aws:utility_lib",
        "//test/extensions/common/aws:aws_mocks",
        "//test/mocks/server:server_factory_context_mocks",
        "//test/test_common:utility_lib",
    ],
)

envoy_cc_test(
    name = "region_provider_impl_test",
    srcs = ["region_provider_impl_test.cc"],
    rbe_pool = "6gig",
    deps = [
        "//source/extensions/common/aws:region_provider_impl_lib",
        "//test/test_common:environment_lib",
    ],
)

envoy_cc_test(
    name = "metadata_fetcher_test",
    srcs = ["metadata_fetcher_test.cc"],
    rbe_pool = "6gig",
    deps = [
        "//source/common/http:headers_lib",
        "//source/common/http:message_lib",
        "//source/common/http:utility_lib",
        "//source/common/protobuf:utility_lib",
        "//source/extensions/common/aws:metadata_fetcher_lib",
        "//test/extensions/common/aws:aws_mocks",
        "//test/extensions/filters/http/common:mock_lib",
        "//test/mocks/api:api_mocks",
        "//test/mocks/event:event_mocks",
        "//test/mocks/server:factory_context_mocks",
        "//test/test_common:environment_lib",
        "//test/test_common:simulated_time_system_lib",
        "//test/test_common:utility_lib",
    ],
)

envoy_cc_test(
    name = "credentials_provider_test",
    srcs = ["credentials_provider_test.cc"],
    rbe_pool = "6gig",
    deps = [
        "//source/extensions/common/aws:credentials_provider_interface",
        "//source/extensions/common/aws:utility_lib",
        "//source/extensions/common/aws/signers:sigv4_signer_impl_lib",
        "//test/extensions/common/aws:aws_mocks",
        "//test/mocks/server:factory_context_mocks",
    ],
)
