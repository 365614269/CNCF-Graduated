#!/bin/bash
# shellcheck disable=SC2034,SC2153,SC2155,SC2164

# Copyright Istio Authors. All Rights Reserved.
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

####################################################################################################
# WARNING: THIS IS AN AUTO-GENERATED FILE, DO NOT EDIT. PLEASE MODIFY THE ORIGINAL MARKDOWN FILE:
#          docs/tasks/traffic-management/ingress/gateway-api/index.md
####################################################################################################

snip_setup_1() {
kubectl get crd gateways.gateway.networking.k8s.io &> /dev/null || \
  { kubectl kustomize "github.com/kubernetes-sigs/gateway-api/config/crd?ref=v1.3.0" | kubectl apply -f -; }
}

snip_setup_2() {
istioctl install --set profile=minimal -y
}

snip_configuring_a_gateway_1() {
kubectl apply -f samples/httpbin/httpbin.yaml
}

snip_configuring_a_gateway_2() {
kubectl create namespace istio-ingress
kubectl apply -f - <<EOF
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: gateway
  namespace: istio-ingress
spec:
  gatewayClassName: istio
  listeners:
  - name: default
    hostname: "*.example.com"
    port: 80
    protocol: HTTP
    allowedRoutes:
      namespaces:
        from: All
---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: http
  namespace: default
spec:
  parentRefs:
  - name: gateway
    namespace: istio-ingress
  hostnames: ["httpbin.example.com"]
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /get
    backendRefs:
    - name: httpbin
      port: 8000
EOF
}

snip_configuring_a_gateway_3() {
kubectl wait -n istio-ingress --for=condition=programmed gateways.gateway.networking.k8s.io gateway
export INGRESS_HOST=$(kubectl get gateways.gateway.networking.k8s.io gateway -n istio-ingress -ojsonpath='{.status.addresses[0].value}')
}

snip_configuring_a_gateway_4() {
curl -s -I -HHost:httpbin.example.com "http://$INGRESS_HOST/get"
}

! IFS=$'\n' read -r -d '' snip_configuring_a_gateway_4_out <<\ENDSNIP
...
HTTP/1.1 200 OK
...
server: istio-envoy
...
ENDSNIP

snip_configuring_a_gateway_5() {
curl -s -I -HHost:httpbin.example.com "http://$INGRESS_HOST/headers"
}

! IFS=$'\n' read -r -d '' snip_configuring_a_gateway_5_out <<\ENDSNIP
HTTP/1.1 404 Not Found
...
ENDSNIP

snip_configuring_a_gateway_6() {
kubectl apply -f - <<EOF
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: http
  namespace: default
spec:
  parentRefs:
  - name: gateway
    namespace: istio-ingress
  hostnames: ["httpbin.example.com"]
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /get
    - path:
        type: PathPrefix
        value: /headers
    filters:
    - type: RequestHeaderModifier
      requestHeaderModifier:
        add:
        - name: my-added-header
          value: added-value
    backendRefs:
    - name: httpbin
      port: 8000
EOF
}

snip_configuring_a_gateway_7() {
curl -s -HHost:httpbin.example.com "http://$INGRESS_HOST/headers" | jq '.headers["My-Added-Header"][0]'
}

! IFS=$'\n' read -r -d '' snip_configuring_a_gateway_7_out <<\ENDSNIP
...
"added-value"
...
ENDSNIP

! IFS=$'\n' read -r -d '' snip_automated_deployment_1 <<\ENDSNIP
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: gateway
spec:
  infrastructure:
    annotations:
      some-key: some-value
    labels:
      key: value
    parametersRef:
      group: ""
      kind: ConfigMap
      name: gw-options
  gatewayClassName: istio
ENDSNIP

! IFS=$'\n' read -r -d '' snip_automated_deployment_2 <<\ENDSNIP
apiVersion: v1
kind: ConfigMap
metadata:
  name: gw-options
data:
  horizontalPodAutoscaler: |
    spec:
      minReplicas: 2
      maxReplicas: 2

  deployment: |
    metadata:
      annotations:
        additional-annotation: some-value
    spec:
      replicas: 4
      template:
        spec:
          containers:
          - name: istio-proxy
            resources:
              requests:
                cpu: 1234m

  service: |
    spec:
      ports:
      - "\$patch": delete
        port: 15021
ENDSNIP

! IFS=$'\n' read -r -d '' snip_gatewayclass_defaults_1 <<\ENDSNIP
kind: IstioOperator
spec:
  values:
    gatewayClasses:
      istio:
        deployment:
          spec:
            replicas: 2
ENDSNIP

! IFS=$'\n' read -r -d '' snip_resource_attachment_and_scaling_1 <<\ENDSNIP
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: gateway
spec:
  gatewayClassName: istio
  listeners:
  - name: default
    hostname: "*.example.com"
    port: 80
    protocol: HTTP
    allowedRoutes:
      namespaces:
        from: All
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: gateway
spec:
  # Match the generated Deployment by reference
  # Note: Do not use `kind: Gateway`.
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: gateway-istio
  minReplicas: 2
  maxReplicas: 5
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 50
---
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: gateway
spec:
  minAvailable: 1
  selector:
    # Match the generated Deployment by label
    matchLabels:
      gateway.networking.k8s.io/gateway-name: gateway
ENDSNIP

! IFS=$'\n' read -r -d '' snip_manual_deployment_1 <<\ENDSNIP
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: gateway
spec:
  addresses:
  - value: ingress.istio-gateways.svc.cluster.local
    type: Hostname
...
ENDSNIP

! IFS=$'\n' read -r -d '' snip_mesh_traffic_1 <<\ENDSNIP
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: mesh
spec:
  parentRefs:
  - group: ""
    kind: Service
    name: example
  rules:
  - filters:
    - type: RequestHeaderModifier
      requestHeaderModifier:
        add:
        - name: my-added-header
          value: added-value
    backendRefs:
    - name: example
      port: 80
ENDSNIP

snip_cleanup_1() {
kubectl delete -f samples/httpbin/httpbin.yaml
kubectl delete httproute http
kubectl delete gateways.gateway.networking.k8s.io gateway -n istio-ingress
kubectl delete ns istio-ingress
}

snip_cleanup_2() {
istioctl uninstall -y --purge
kubectl delete ns istio-system
}

snip_cleanup_3() {
kubectl kustomize "github.com/kubernetes-sigs/gateway-api/config/crd?ref=v1.3.0" | kubectl delete -f -
}
