#!/usr/bin/env bash
# shellcheck disable=SC1090,SC2154

# Copyright 2020 Istio Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# @setup profile=demo

set -e
set -u
set -o pipefail

GATEWAY_API="${GATEWAY_API:-false}"

# Make sure automatic sidecar injection is enabled
kubectl label namespace default istio-injection=enabled || true

# Deploy curl sample
snip_before_you_begin_1
_wait_for_deployment default curl

# Generate Certificates for service outside the mesh to use for mTLS
set +e # suppress harmless "No such file or directory:../crypto/bio/bss_file.c:72:fopen('1_root/index.txt.attr','r')" error
snip_generate_client_and_server_certificates_and_keys_1
snip_generate_client_and_server_certificates_and_keys_2
snip_generate_client_and_server_certificates_and_keys_4
set -e

# Create mesh-external namespace
snip_deploy_a_mutual_tls_server_1

# Setup sever with certs and config
snip_deploy_a_mutual_tls_server_2
snip_deploy_a_mutual_tls_server_3
snip_deploy_a_mutual_tls_server_4
snip_deploy_a_mutual_tls_server_5

# Wait for nginx
_wait_for_deployment mesh-external my-nginx

# Open Gateway Listener
if [ "$GATEWAY_API" == "true" ]; then
    snip_configure_mutual_tls_origination_for_egress_traffic_2
    snip_configure_mutual_tls_origination_for_egress_traffic_4
else
    snip_configure_mutual_tls_origination_for_egress_traffic_1
    snip_configure_mutual_tls_origination_for_egress_traffic_3
    _wait_for_resource gateway default istio-egressgateway
    _wait_for_resource destinationrule default egressgateway-for-nginx
fi

# Configure routing from curl to egress gateway to nginx
if [ "$GATEWAY_API" == "true" ]; then
    snip_configure_mutual_tls_origination_for_egress_traffic_6
else
    snip_configure_mutual_tls_origination_for_egress_traffic_5
    _wait_for_resource virtualservice default direct-nginx-through-egress-gateway
fi

# Originate TLS with destination rule
if [ "$GATEWAY_API" == "true" ]; then
    snip_configure_mutual_tls_origination_for_egress_traffic_8
    _verify_contains snip_configure_mutual_tls_origination_for_egress_traffic_10 "kubernetes://client-credential            Cert Chain     ACTIVE"
else
    snip_configure_mutual_tls_origination_for_egress_traffic_7
    _wait_for_resource destinationrule istio-system originate-mtls-for-nginx

    _verify_contains snip_configure_mutual_tls_origination_for_egress_traffic_9 "kubernetes://client-credential            Cert Chain     ACTIVE"
fi

# Verify that mTLS connection is set up properly
_verify_contains snip_configure_mutual_tls_origination_for_egress_traffic_11 "Welcome to nginx!"

# Verify request is routed through Gateway
if [ "$GATEWAY_API" == "true" ]; then
    _verify_contains snip_configure_mutual_tls_origination_for_egress_traffic_13 "GET / HTTP/1.1"
else
    _verify_contains snip_configure_mutual_tls_origination_for_egress_traffic_12 "GET / HTTP/1.1"
fi

# @cleanup
if [ "$GATEWAY_API" != "true" ]; then
    kubectl label namespace default istio-injection-
    snip_cleanup_the_mutual_tls_origination_example_1
    snip_cleanup_the_mutual_tls_origination_example_2
    snip_cleanup_the_mutual_tls_origination_example_4
    snip_cleanup_the_mutual_tls_origination_example_5
    snip_cleanup_1
fi
