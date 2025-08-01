---
title: Install Primary-Remote on different networks
description: Install an Istio mesh across primary and remote clusters on different networks.
weight: 40
keywords: [kubernetes,multicluster]
test: yes
owner: istio/wg-environments-maintainers
---
Follow this guide to install the Istio control plane on `cluster1` (the
{{< gloss >}}primary cluster{{< /gloss >}}) and configure `cluster2` (the
{{< gloss >}}remote cluster{{< /gloss >}}) to use the control plane in
`cluster1`. Cluster `cluster1` is on the `network1` network, while `cluster2`
is on the `network2` network. This means there is no direct connectivity
between pods across cluster boundaries.

Before proceeding, be sure to complete the steps under
[before you begin](/es/docs/setup/install/multicluster/before-you-begin).

{{< boilerplate multi-cluster-with-metallb >}}

In this configuration, cluster `cluster1` will observe the API Servers in
both clusters for endpoints. In this way, the control plane will be able to
provide service discovery for workloads in both clusters.

Service workloads across cluster boundaries communicate indirectly, via
dedicated gateways for [east-west](https://en.wikipedia.org/wiki/East-west_traffic)
traffic. The gateway in each cluster must be reachable from the other cluster.

Services in `cluster2` will reach the control plane in `cluster1` via the
same east-west gateway.

{{< image width="75%"
    link="arch.svg"
    caption="Primary and remote clusters on separate networks"
    >}}

## Set the default network for `cluster1`

If the istio-system namespace is already created, we need to set the cluster's network there:

{{< text bash >}}
$ kubectl --context="${CTX_CLUSTER1}" get namespace istio-system && \
  kubectl --context="${CTX_CLUSTER1}" label namespace istio-system topology.istio.io/network=network1
{{< /text >}}

## Configure `cluster1` as a primary

Create the `istioctl` configuration for `cluster1`:

{{< tabset category-name="multicluster-primary-remote-install-type-primary-cluster" >}}

{{< tab name="IstioOperator" category-value="iop" >}}

Install Istio as primary in `cluster1` using istioctl and the `IstioOperator` API.

{{< text bash >}}
$ cat <<EOF > cluster1.yaml
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
spec:
  values:
    global:
      meshID: mesh1
      multiCluster:
        clusterName: cluster1
      network: network1
      externalIstiod: true
EOF
{{< /text >}}

Apply the configuration to `cluster1`:

{{< text bash >}}
$ istioctl install --context="${CTX_CLUSTER1}" -f cluster1.yaml
{{< /text >}}

Notice that `values.global.externalIstiod` is set to `true`. This enables the control plane
installed on `cluster1` to also serve as an external control plane for other remote clusters.
When this feature is enabled, `istiod` will attempt to acquire the leadership lock, and consequently manage,
[appropriately annotated](#set-the-control-plane-cluster-for-cluster2) remote clusters that are
attached to it (`cluster2` in this case).

{{< /tab >}}

{{< tab name="Helm" category-value="helm" >}}

Install Istio as primary in `cluster1` using the following Helm commands:

Install the `base` chart in `cluster1`:

{{< text bash >}}
$ helm install istio-base istio/base -n istio-system --kube-context "${CTX_CLUSTER1}"
{{< /text >}}

Then, install the `istiod` chart in `cluster1` with the following multi-cluster settings:

{{< text bash >}}
$ helm install istiod istio/istiod -n istio-system --kube-context "${CTX_CLUSTER1}" --set global.meshID=mesh1 --set global.externalIstiod=true --set global.multiCluster.clusterName=cluster1 --set global.network=network1
{{< /text >}}

Notice that `values.global.externalIstiod` is set to `true`. This enables the control plane
installed on `cluster1` to also serve as an external control plane for other remote clusters.
When this feature is enabled, `istiod` will attempt to acquire the leadership lock, and consequently manage,
[appropriately annotated](#set-the-control-plane-cluster-for-cluster2) remote clusters that are
attached to it (`cluster2` in this case).

{{< /tab >}}

{{< /tabset >}}

## Install the east-west gateway in `cluster1`

Install a gateway in `cluster1` that is dedicated to east-west traffic. By
default, this gateway will be public on the Internet. Production systems may
require additional access restrictions (e.g. via firewall rules) to prevent
external attacks. Check with your cloud vendor to see what options are
available.

{{< tabset category-name="east-west-gateway-install-type-cluster-1" >}}

{{< tab name="IstioOperator" category-value="iop" >}}

{{< text bash >}}
$ @samples/multicluster/gen-eastwest-gateway.sh@ \
    --network network1 | \
    istioctl --context="${CTX_CLUSTER1}" install -y -f -
{{< /text >}}

{{< warning >}}
If the control-plane was installed with a revision, add the `--revision rev` flag to the `gen-eastwest-gateway.sh` command.
{{< /warning >}}

{{< /tab >}}
{{< tab name="Helm" category-value="helm" >}}

Install the east-west gateway in `cluster1` using the following Helm command:

{{< text bash >}}
$ helm install istio-eastwestgateway istio/gateway -n istio-system --kube-context "${CTX_CLUSTER1}" --set name=istio-eastwestgateway --set networkGateway=network1
{{< /text >}}

{{< warning >}}
If the control-plane was installed with a revision, you must add a `--set revision=<my-revision>` flag to the Helm install command.
{{< /warning >}}

{{< /tab >}}

{{< /tabset >}}

Wait for the east-west gateway to be assigned an external IP address:

{{< text bash >}}
$ kubectl --context="${CTX_CLUSTER1}" get svc istio-eastwestgateway -n istio-system
NAME                    TYPE           CLUSTER-IP    EXTERNAL-IP    PORT(S)   AGE
istio-eastwestgateway   LoadBalancer   10.80.6.124   34.75.71.237   ...       51s
{{< /text >}}

## Expose the control plane in `cluster1`

Before we can install on `cluster2`, we need to first expose the control plane in
`cluster1` so that services in `cluster2` will be able to access service discovery:

{{< text bash >}}
$ kubectl apply --context="${CTX_CLUSTER1}" -n istio-system -f \
    @samples/multicluster/expose-istiod.yaml@
{{< /text >}}

{{< warning >}}
If the control-plane was installed with a revision `rev`, use the following command instead:

{{< text bash >}}
$ sed 's/{{.Revision}}/rev/g' @samples/multicluster/expose-istiod-rev.yaml.tmpl@ | kubectl apply --context="${CTX_CLUSTER1}" -n istio-system -f -
{{< /text >}}

{{< /warning >}}

## Set the control plane cluster for `cluster2`

We need identify the external control plane cluster that should manage `cluster2` by annotating the
istio-system namespace:

{{< text bash >}}
$ kubectl --context="${CTX_CLUSTER2}" create namespace istio-system
$ kubectl --context="${CTX_CLUSTER2}" annotate namespace istio-system topology.istio.io/controlPlaneClusters=cluster1
{{< /text >}}

Setting the `topology.istio.io/controlPlaneClusters` namespace annotation to `cluster1` instructs the `istiod`
running in the same namespace (istio-system in this case) on `cluster1` to manage `cluster2` when it
is [attached as a remote cluster](#attach-cluster2-as-a-remote-cluster-of-cluster1).

## Set the default network for `cluster2`

Set the network for `cluster2` by adding a label to the istio-system namespace:

{{< text bash >}}
$ kubectl --context="${CTX_CLUSTER2}" label namespace istio-system topology.istio.io/network=network2
{{< /text >}}

## Configure `cluster2` as a remote

Save the address of `cluster1`’s east-west gateway.

{{< text bash >}}
$ export DISCOVERY_ADDRESS=$(kubectl \
    --context="${CTX_CLUSTER1}" \
    -n istio-system get svc istio-eastwestgateway \
    -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
{{< /text >}}

Now create a remote configuration on `cluster2`.

{{< tabset category-name="multicluster-primary-remote-install-type-remote-cluster" >}}

{{< tab name="IstioOperator" category-value="iop" >}}

{{< text bash >}}
$ cat <<EOF > cluster2.yaml
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
spec:
  profile: remote
  values:
    istiodRemote:
      injectionPath: /inject/cluster/cluster2/net/network2
    global:
      remotePilotAddress: ${DISCOVERY_ADDRESS}
EOF
{{< /text >}}

Apply the configuration to `cluster2`:

{{< text bash >}}
$ istioctl install --context="${CTX_CLUSTER2}" -f cluster2.yaml
{{< /text >}}

{{< /tab >}}
{{< tab name="Helm" category-value="helm" >}}

Install Istio as remote in `cluster2` using the following Helm commands:

Install the `base` chart in `cluster2`:

{{< text bash >}}
$ helm install istio-base istio/base -n istio-system --set profile=remote --kube-context "${CTX_CLUSTER2}"
{{< /text >}}

Then, install the `istiod` chart in `cluster2` with the following multi-cluster settings:

{{< text bash >}}
$ helm install istiod istio/istiod -n istio-system --set profile=remote --set global.multiCluster.clusterName=cluster2 --set global.network=network2 --set istiodRemote.injectionPath=/inject/cluster/cluster2/net/network2  --set global.configCluster=true --set global.remotePilotAddress="${DISCOVERY_ADDRESS}" --kube-context "${CTX_CLUSTER2}"
{{< /text >}}

{{< tip >}}

The `remote` profile for the `base` and `istiod` Helm charts is only available from Istio release 1.24 onwards.

{{< /tip >}}

{{< /tab >}}

{{< /tabset >}}

{{< tip >}}
Here we're configuring the location of the control plane using the `injectionPath` and
`remotePilotAddress` parameters. Although convenient for demonstration, in a production
environment it is recommended to instead configure the `injectionURL` parameter using
properly signed DNS certs similar to the configuration shown in the
[external control plane instructions](/es/docs/setup/install/external-controlplane/#register-the-new-cluster).
{{< /tip >}}

## Attach `cluster2` as a remote cluster of `cluster1`

To attach the remote cluster to its control plane, we give the control
plane in `cluster1` access to the API Server in `cluster2`. This will do the
following:

- Enables the control plane to authenticate connection requests from
  workloads running in `cluster2`. Without API Server access, the control
  plane will reject the requests.

- Enables discovery of service endpoints running in `cluster2`.

Because it has been included in the `topology.istio.io/controlPlaneClusters` namespace
annotation, the control plane on `cluster1` will also:

- Patch certs in the webhooks in `cluster2`.

- Start the namespace controller which writes configmaps in namespaces in `cluster2`.

To provide API Server access to `cluster2`, we generate a remote secret and
apply it to `cluster1`:

{{< text bash >}}
$ istioctl create-remote-secret \
    --context="${CTX_CLUSTER2}" \
    --name=cluster2 | \
    kubectl apply -f - --context="${CTX_CLUSTER1}"
{{< /text >}}

## Install the east-west gateway in `cluster2`

As we did with `cluster1` above, install a gateway in `cluster2` that is dedicated
to east-west traffic and expose user services.

{{< tabset category-name="east-west-gateway-install-type-cluster-2" >}}

{{< tab name="IstioOperator" category-value="iop" >}}

{{< text bash >}}
$ @samples/multicluster/gen-eastwest-gateway.sh@ \
    --network network2 | \
    istioctl --context="${CTX_CLUSTER2}" install -y -f -
{{< /text >}}

{{< /tab >}}
{{< tab name="Helm" category-value="helm" >}}

Install the east-west gateway in `cluster2` using the following Helm command:

{{< text bash >}}
$ helm install istio-eastwestgateway istio/gateway -n istio-system --kube-context "${CTX_CLUSTER2}" --set name=istio-eastwestgateway --set networkGateway=network2
{{< /text >}}

{{< warning >}}
If the control-plane was installed with a revision, you must add a `--set revision=<my-revision>` to the Helm install command.
{{< /warning >}}

{{< /tab >}}

{{< /tabset >}}

Wait for the east-west gateway to be assigned an external IP address:

{{< text bash >}}
$ kubectl --context="${CTX_CLUSTER2}" get svc istio-eastwestgateway -n istio-system
NAME                    TYPE           CLUSTER-IP    EXTERNAL-IP    PORT(S)   AGE
istio-eastwestgateway   LoadBalancer   10.0.12.121   34.122.91.98   ...       51s
{{< /text >}}

## Expose services in `cluster1` and `cluster2`

Since the clusters are on separate networks, we also need to expose all user
services (*.local) on the east-west gateway in both clusters. While these
gateways are public on the Internet, services behind them can only be accessed by
services with a trusted mTLS certificate and workload ID, just as if they were
on the same network.

{{< text bash >}}
$ kubectl --context="${CTX_CLUSTER1}" apply -n istio-system -f \
    @samples/multicluster/expose-services.yaml@
{{< /text >}}

{{< tip >}}
Since `cluster2` is installed with a remote profile, exposing services on the primary cluster will expose them on the east-west gateways of both clusters.
{{< /tip >}}

**Congratulations!** You successfully installed an Istio mesh across primary
and remote clusters on different networks!

## Next Steps

You can now [verify the installation](/es/docs/setup/install/multicluster/verify).

## Cleanup

Uninstall Istio from both `cluster1` and `cluster2` using the same mechanism you installed Istio with (istioctl or Helm).

{{< tabset category-name="multicluster-uninstall-type-cluster-1" >}}

{{< tab name="IstioOperator" category-value="iop" >}}

Uninstall Istio in `cluster1`:

{{< text syntax=bash snip_id=none >}}
$ istioctl uninstall --context="${CTX_CLUSTER1}" -y --purge
$ kubectl delete ns istio-system --context="${CTX_CLUSTER1}"
{{< /text >}}

Uninstall Istio in `cluster2`:

{{< text syntax=bash snip_id=none >}}
$ istioctl uninstall --context="${CTX_CLUSTER2}" -y --purge
$ kubectl delete ns istio-system --context="${CTX_CLUSTER2}"
{{< /text >}}

{{< /tab >}}

{{< tab name="Helm" category-value="helm" >}}

Delete Istio Helm installation from `cluster1`:

{{< text syntax=bash >}}
$ helm delete istiod -n istio-system --kube-context "${CTX_CLUSTER1}"
$ helm delete istio-eastwestgateway -n istio-system --kube-context "${CTX_CLUSTER1}"
$ helm delete istio-base -n istio-system --kube-context "${CTX_CLUSTER1}"
{{< /text >}}

Delete the `istio-system` namespace from `cluster1`:

{{< text syntax=bash >}}
$ kubectl delete ns istio-system --context="${CTX_CLUSTER1}"
{{< /text >}}

Delete Istio Helm installation from `cluster2`:

{{< text syntax=bash >}}
$ helm delete istiod -n istio-system --kube-context "${CTX_CLUSTER2}"
$ helm delete istio-eastwestgateway -n istio-system --kube-context "${CTX_CLUSTER2}"
$ helm delete istio-base -n istio-system --kube-context "${CTX_CLUSTER2}"
{{< /text >}}

Delete the `istio-system` namespace from `cluster2`:

{{< text syntax=bash >}}
$ kubectl delete ns istio-system --context="${CTX_CLUSTER2}"
{{< /text >}}

(Optional) Delete CRDs installed by Istio:

Deleting CRDs permanently removes any Istio resources you have created in your clusters.
To delete Istio CRDs installed in your clusters:

{{< text syntax=bash snip_id=delete_crds >}}
$ kubectl get crd -oname --context "${CTX_CLUSTER1}" | grep --color=never 'istio.io' | xargs kubectl delete --context "${CTX_CLUSTER1}"
$ kubectl get crd -oname --context "${CTX_CLUSTER2}" | grep --color=never 'istio.io' | xargs kubectl delete --context "${CTX_CLUSTER2}"
{{< /text >}}

{{< /tab >}}

{{< /tabset >}}
