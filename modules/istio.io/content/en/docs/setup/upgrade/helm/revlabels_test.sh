#!/usr/bin/env bash
# Copyright Istio Authors
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
source "content/en/docs/setup/upgrade/helm/common.sh"
source "content/en/boilerplates/snips/crd-upgrade-123.sh"

set -e
set -u

set -o pipefail

# @setup profile=none
_install_istio_helm

_rewrite_helm_repo snip_usage_1

_rewrite_helm_repo snip_usage_2
_rewrite_helm_repo snip_default_tag_1

_remove_istio_helm

kubectl delete mutatingwebhookconfiguration istio-revision-tag-default
kubectl delete mutatingwebhookconfiguration istio-revision-tag-prod-canary
kubectl delete mutatingwebhookconfiguration istio-revision-tag-prod-stable

# @cleanup
