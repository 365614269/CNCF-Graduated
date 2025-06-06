---
title: Installing Linkerd with Helm
description: Install Linkerd onto your Kubernetes cluster using Helm.
---

Linkerd can be installed via Helm rather than with the `linkerd install`
command. This is recommended for production, since it allows for repeatability.

{{< docs/production-note >}}

## Prerequisite: generate identity certificates

To do [automatic mutual TLS](../features/automatic-mtls/), Linkerd requires
trust anchor certificate and an issuer certificate and key pair. When you're
using `linkerd install`, we can generate these for you. However, for Helm,
you will need to generate these yourself.

Please follow the instructions in [Generating your own mTLS root
certificates](generate-certificates/) to generate these.

## Helm install procedure for stable releases

```bash
# add the repo for stable releases:
helm repo add linkerd https://helm.linkerd.io/stable

helm install linkerd2 \
  --set-file identityTrustAnchorsPEM=ca.crt \
  --set-file identity.issuer.tls.crtPEM=issuer.crt \
  --set-file identity.issuer.tls.keyPEM=issuer.key \
  linkerd/linkerd2
```

The chart values will be picked from the chart's `values.yaml` file.

You can override the values in that file by providing your own `values.yaml`
file passed with a `-f` option, or overriding specific values using the family of
`--set` flags like we did above for certificates.

## Helm install procedure for edge releases

You need to install two separate charts in succession: first `linkerd-crds` and
then `linkerd-control-plane`. This new method will eventually make its way into
the `2.12.0` stable release as well when it comes out.

{{< note >}}
If installing Linkerd in a cluster that uses Cilium in kube-proxy replacement
mode, additional steps may be needed to ensure service discovery works as
intended. Instrunctions are on the [Cilium cluster
configuration](../reference/cluster-configuration/#cilium) page.
{{< /note >}}

### linkerd-crds

The `linkerd-crds` chart sets up the CRDs linkerd requires:

```bash
# add the repo for edge releases:
helm repo add linkerd-edge https://helm.linkerd.io/edge

helm install linkerd-crds -n linkerd --create-namespace --devel linkerd-edge/linkerd-crds
```

{{< note >}}
This will create the `linkerd` namespace. If it already exists or you're
creating it beforehand elsewhere in your pipeline, just omit the
`--create-namespace` flag.
{{< /note >}}

### linkerd-control-plane

The `linkerd-control-plane` chart sets up all the control plane components:

```bash
helm install linkerd-control-plane \
  -n linkerd \
  --devel \
  --set-file identityTrustAnchorsPEM=ca.crt \
  --set-file identity.issuer.tls.crtPEM=issuer.crt \
  --set-file identity.issuer.tls.keyPEM=issuer.key \
  linkerd-edge/linkerd-control-plane
```

{{< note >}}
If you are using [Linkerd's CNI plugin](../features/cni/), you must also add the
`--set cniEnabled=true` flag to your `helm install` command.
{{< /note >}}

## Enabling high availability mode

The linkerd2 chart (or, for edge releases, the linkerd-control-plane chart)
contains a file called `values-ha.yaml` that overrides some default values to
enable high availability mode, analogous to the `--ha` option in `linkerd
install`.

You can get the `values-ha.yaml` by fetching the chart files:

```bash
# for stable
helm fetch --untar linkerd/linkerd2

# for edge
helm fetch --untar --devel linkerd-edge/linkerd-control-plane
```

Then use the `-f` flag to provide this override file. For example:

```bash
# for stable
helm install linkerd2 \
  --set-file identityTrustAnchorsPEM=ca.crt \
  --set-file identity.issuer.tls.crtPEM=issuer.crt \
  --set-file identity.issuer.tls.keyPEM=issuer.key \
  -f linkerd2/values-ha.yaml \
  linkerd/linkerd2

# for edge
helm install linkerd-control-plane \
  -n linkerd \
  --devel \
  --set-file identityTrustAnchorsPEM=ca.crt \
  --set-file identity.issuer.tls.crtPEM=issuer.crt \
  --set-file identity.issuer.tls.keyPEM=issuer.key \
  -f linkerd-control-plane/values-ha.yaml \
  linkerd-edge/linkerd-control-plane
```

## Customizing the namespace

To install Linkerd to a different namespace, you can override the Helm
`Namespace` variable.

By default, the chart creates the control plane namespace with the
`config.linkerd.io/admission-webhooks: disabled` label. This is required for the
control plane to work correctly. This means that the chart won't work with
Helm's `--namespace` option.  If you're relying on a separate tool to create the
control plane namespace, make sure that:

1. The namespace is labeled with `config.linkerd.io/admission-webhooks: disabled`
1. The `installNamespace` is set to `false`
1. The `namespace` variable is overridden with the name of your namespace

## Upgrading with Helm

First, make sure your local Helm repos are updated:

```bash
helm repo update

helm search repo linkerd2
NAME                    CHART VERSION          APP VERSION            DESCRIPTION
linkerd/linkerd2        <chart-semver-version> {{< latest-stable-version >}}    Linkerd gives you observability, reliability, and securit...
```

During an upgrade, you must choose whether you want to reuse the values in the
chart or move to the values specified in the newer chart.  Our advice is to use
a `values.yaml` file that stores all custom overrides that you have for your
chart.

The `helm upgrade` command has a number of flags that allow you to customize its
behavior. Special attention should be paid to `--reuse-values` and
`--reset-values` and how they behave when charts change from version to version
and/or overrides are applied through `--set` and `--set-file`.  For example:

- `--reuse-values` with no overrides - all values are reused
- `--reuse-values` with overrides - all except the values that are overridden
are reused
- `--reset-values` with no overrides - no values are reused and all changes
from provided release are applied during the upgrade
- `--reset-values` with overrides - no values are reused and changed from
provided release are applied together with the overrides
- no flag and no overrides - `--reuse-values` will be used by default
- no flag and overrides - `--reset-values` will be used by default

Finally, before upgrading, check whether there are breaking changes to the chart
(i.e. renamed or moved keys, etc). You can consult the
[edge](https://hub.helm.sh/charts/linkerd2-edge/linkerd2) or the
[stable](https://hub.helm.sh/charts/linkerd2/linkerd2) chart docs, depending on
which one your are upgrading to. If there are, make the corresponding changes to
your `values.yaml` file. Then you can use:

```bash
helm upgrade linkerd2 linkerd/linkerd2 --reset-values -f values.yaml --atomic
```

The `--atomic` flag will ensure that all changes are rolled back in case the
upgrade operation fails.
