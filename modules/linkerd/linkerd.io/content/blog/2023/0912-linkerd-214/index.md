---
date: 2023-09-12T00:00:00Z
slug: linkerd-214
title: |-
  Workshop Recap: A closer look at flat-network multicluster and HTTPRoute timeouts with Linkerd 2.14
keywords: [linkerd, "2.14", features]
params:
  author: flynn
  showCover: true
---

_This blog post is based on a workshop that I recently delivered at Buoyant’s
[Service Mesh Academy](https://buoyant.io/service-mesh-academy). If this seems
interesting, check out the
[full recording](https://buoyant.io/service-mesh-academy/whats-coming-in-linkerd-2-14/)!_

Linkerd 2.14 introduces several new features, including flat-network
multicluster, workload identity across clusters, Gateway API conformance,
timeout support in HTTPRoute, and (of course) many bugfixes. Let's take a closer
look at some of the highlights.

## Flat-Network Multicluster

In Linkerd 2.13 and earlier, Linkerd's multicluster functionality worked by
channeling traffic through multicluster gateways.

![Linkerd 2.13 Gateways](gateways.png "Linkerd 2.13 multicluster gateways")

This mechanism is simple and functions very well as long as the gateways have IP
connectivity, but it has two significant caveats:

1. **Added Cost:** Funnelling traffic through the gateways adds latency (though
   not very much in Linkerd's case). Additionally, the gateways themselves often
   require network load balancers, which can add expense.

2. **Lost Identity:** When a workload in one cluster calls a workload in a
   different cluster, the workload in the second cluster sees the request
   originating from its gateway, not from the workload in the first cluster.
   This makes it impossible to properly implement authorization policy across
   clusters.

![Linkerd 2.13 Lost Identity](lost-identity.png "Linkerd 2.13 multicluster lost identity")

The solution in 2.14? _Completely bypass the gateways_.

![Linkerd 2.14 Flat-network Multicluster](flat-multicluster.png "Linkerd 2.14 flat-network multicluster")

This new approach brings some real benefits to the table, all stemming from the
fact that now workload 1 is simply making a normal mTLS connection directly to
workload 2. This means that workload identity _is_ preserved across cluster
boundaries, since the request to workload 2 comes directly from workload 1. In
turn, this means that authorization policy just works across cluster boundaries,
as does all the rest of Linkerd's functionality related to mTLS and identity.

Additionally, this direct pod-to-pod communication reduces latency and also
means that there are fewer components to fail in the network path. There's also
no need for load balancers, which can directly cut costs.

To use this new capability, you start by installing the Linkerd multicluster
extension with gateways disabled:

```bash
linkerd multicluster install --gateway=false | kubectl apply -f -
```

and then also add `--gateway=false` when linking clusters, e.g.:

```bash
linkerd --context=us-west multicluster link \
        --cluster-name us-west \
        --gateway=false \
    | kubectl --context=us-east apply -f -
```

Note that even though the gateways are no longer in play, you'll still see the
Link resources.

Finally, when you label Services for export across clusters, you'll use

```yaml
mirror.linkerd.io/exported: remote-discovery
```

The `remote-discovery` value for this label triggers Linkerd to not use
gateways, and instead directly communicate between the two clusters' control
planes to handle Service mirroring across clusters.

Of course, there are still points to be aware of. The most critical is also the
most obvious: you must have a flat network in which each cluster has its own
unique CIDR range, and Pods in one cluster can talk directly to Pods in other
clusters. How tricky this is to arrange will depend on your cluster provider,
but it's always possible!

Additionally - as always with multicluster operations in Linkerd - your clusters
will all need a shared trust anchor. We generally recommend that each cluster
have a distinct identity trust domain (and thus cluster domain), as this can
considerably simplify cross-cluster authorization policy.

### Gateway API Improvements

Starting in Linkerd 2.12, the Gateway API became the central configuration
mechanism for talking about classes of HTTP traffic (including gRPC), including
functionality such as authentication policy (Linkerd 2.12) and dynamic request
routing (Linkerd 2.13). In the long run, Gateway API is expected to replace SMI
and ServicePolicy, though there's no defined timetable for that yet.

A challenge when implementing Gateway API as a service mesh is that the
conformance tests for Gateway API were originally defined assuming that any
implementation being tested for conformance included a gateway controller that
handled north/south traffic. Since Linkerd relies on external ingress
controllers for north/south traffic rather than bundling its own, there was no
way that Linkerd could be conformant.

This situation changes in Gateway API v0.8.0, just released on August 29, 2023,
which adds the concept of _conformance profiles_ to Gateway API. These named
profiles define subsets of Gateway API and allow implementations to choose (and
document) the subsets to which they conform. One of the profiles defined by
v0.8.0 is the `Mesh` profile, which checks only service mesh functionality, and
we're delighted to announce that Linkerd 2.14 is fully conformant with this new
`Mesh` profile.

A major component of Gateway API conformance is, of course, the APIGroup used
for Gateway API resources. Since Linkerd 2.13 and earlier could not be
conformant, by definition, they copied the HTTPRoute resource into the
`policy.linkerd.io` APIGroup. Linkerd 2.14 adds support for the official Gateway
API group, `gateway.networking.k8s.io`. New installations should generally use
`gateway.networking.k8s.io` HTTPRoutes (with one notable exception, described
below), but of course support for `policy.linkerd.io` will remain for the
foreseeable future.

## HTTPRoute Timeouts

One more thing that Linkerd 2.14 adds is support for timeouts in HTTPRoutes:

```yaml
apiVersion: policy.linkerd.io/v1beta3
kind: HTTPRoute
metadata:
  name: timeout-example
spec:
  ...
  rules:
  - backendRefs:
    - name: some-service
      port: 8080
    timeouts:
      request: 10s
      backendRequest: 3s
```

The two types of timeouts have different purposes:

- `timeouts.request` sets the end-to-end timeout for the entire request,
  including retries, etc., and
- `timeouts.backendRequest` sets the timeout for just the communications from
  the proxies to the workload.

<!-- ![Linkerd 2.14 HTTPRoute Timeouts](httproute-timeouts.png "Linkerd 2.14 HTTPRoute timeout semantics") -->

Either timeout can be set alone, or both can be set together (in which case
`timeouts.backendRequest` must be less than `timeouts.request`). To disable a
timeout, set it to `0s` -- _not_ setting a timeout leaves Linkerd to enforce its
default.

Both timeouts use the syntax of Go's `time.ParseDuration`, with the exception
that floating-point numbers are not allowed (so `1h30m` is legal, but `1.5h` is
not). The gory details are specified in GEP-2257 and GEP-1742, but it's worth
noting in particular that the zero value for "no timeout" must be specified as
`0s`, not simply `0`.

There's one final important point about timeouts: Linkerd 2.14 only supports
timeouts in `policy.linkerd.io/v1beta3` HTTPRoutes, because GEP-1742 didn't
quite make it into Gateway API v0.8.0 (it's slated for inclusion in Gateway API
v1.0.0). This is the one reason you might still need to use `policy.linkerd.io`
HTTPRoutes instead of switching to `gateway.networking.k8s.io` HTTPRoutes.

## Linkerd 2.14

Linkerd has a long tradition of easily providing security, reliability, and
observability, while maintaining operational simplicity and performance. Linkerd
2.14 continues this tradition, adding important new capabilities like
flat-network multicluster and enhanced Gateway API functionality in order to
address the needs of users from small startups to huge enterprises.

---

If you found this interesting, check out the Service Mesh Academy workshop on
[What's Coming in Linkerd 2.14](https://buoyant.io/service-mesh-academy/whats-coming-in-linkerd-2-14/),
where you can see the hands-on demo of everything I've talked about here! And,
as always, feedback is always welcome -- you can find me as `@flynn` on the
[Linkerd Slack](https://slack.linkerd.io).
