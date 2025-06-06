---
title: IPTables Reference
description: A table with all of the chains and associated rules
---

In order to route TCP traffic in a pod to and from the proxy, an [`init
container`](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/)
is used to set up `iptables` rules at the start of an injected pod's
lifecycle.

At first, `linkerd-init` will create two chains in the `nat` table:
`PROXY_INIT_REDIRECT`, and `PROXY_INIT_OUTPUT`. These chains are used to route
inbound and outbound packets through the proxy. Each chain has a set of rules
attached to it, these rules are traversed by a packet in order.

## Inbound connections

When a packet arrives in a pod, it will typically be processed by the
`PREROUTING` chain, a default chain attached to the `nat` table. The sidecar
container will create a new chain to process inbound packets, called
`PROXY_INIT_REDIRECT`.  The sidecar container creates a rule
(`install-proxy-init-prerouting`) to send packets from the `PREROUTING` chain
to our redirect chain. This is the first rule traversed by an inbound packet.

The redirect chain will be configured with two more rules:

1. `ignore-port`: will ignore processing packets whose destination ports are
     included in the `skip-inbound-ports` install option.
2. `proxy-init-redirect-all`: will redirect all incoming TCP packets through
     the proxy, on port `4143`.

Based on these two rules, there are two possible paths that an inbound packet
can take, both of which are outlined below.

![Inbound iptables chain traversal](/docs/images/iptables/iptables-fig2-1.png "Inbound iptables chain traversal")

The packet will arrive on the `PREROUTING` chain and will be immediately routed
to the redirect chain. If its destination port matches any of the inbound ports
to skip, then it will be forwarded directly to the application process,
_bypassing the proxy_. The list of destination ports to check against can be
[configured when installing Linkerd](cli/install/#). If the
packet does not match any of the ports in the list, it will be redirected
through the proxy. Redirection is done by changing the incoming packet's
destination header, the target port will be replaced with `4143`, which is the
proxy's inbound port. The proxy will process the packet and produce a new one
that will be forwarded to the service; it will be able to get the original
target (IP:PORT) of the inbound packet by using a special socket option
[`SO_ORIGINAL_DST`](https://linux.die.net/man/3/getsockopt). The new packet
will be routed through the `OUTPUT` chain, from there it will be sent to the
application. The `OUTPUT` chain rules are covered in more detail below.

## Outbound connections

When a packet leaves a pod, it will first traverse the `OUTPUT` chain, the
first default chain an outgoing packet traverses in the `nat` table. To
redirect outgoing packets through the outbound side of the proxy, the sidecar
container will again create a new chain. The first outgoing rule is similar to
the inbound counterpart: any packet that traverses the `OUTPUT` chain should be
forwarded to our `PROXY_INIT_OUTPUT` chain to be processed.

The output redirect chain is slightly harder to understand but follows the same
logical flow as the inbound redirect chain, in total there are 4 rules
configured:

1. `ignore-proxy-uid`: any packets owned by the proxy (whose user id is
     `2102`), will skip processing and return to the previous (`OUTPUT`) chain.
     From there, it will be sent on the outbound network interface (either to
     the application, in the case of an inbound packet, or outside of the pod,
     for an outbound packet).
2. `ignore-loopback`: if the packet is sent over the loopback interface
     (`lo`), it will skip processing and return to the previous chain. From
     here, the packet will be sent to the destination, much like the first rule
     in the chain.
3. `ignore-port`: will ignore processing packets whose destination ports are
     included in the `skip-outbound-ports` install option.
4. `redirect-all-outgoing`: the last rule in the chain, it will redirect all
     outgoing TCP packets to port `4140`, the proxy's outbound port. If a
     packet has made it this far, it is guaranteed its destination is not local
     (i.e `lo`) and it has not been produced by the proxy. This means the
     packet has been produced by the service, so it should be forwarded to its
     destination by the proxy.

![Outbound iptables chain traversal](/docs/images/iptables/iptables-fig2-2.png "Outbound iptables chain traversal")

A packet produced by the service will first hit the `OUTPUT` chain; from here,
it will be sent to our own output chain for processing. The first rule it
encounters in `PROXY_INIT_OUTPUT` will be `ignore-proxy-uid`. Since the packet
was generated by the service, this rule will be skipped. If the packet's
destination is not a port bound on localhost (e.g `127.0.0.1:80`), then it will
skip the second rule as well. The third rule, `ignore-port` will be matched if
the packet's destination port is in the outbound ports to skip list, in this
case, it will be sent out on the network interface, bypassing the proxy. If the
rule is not matched, then the packet will reach the final rule in the chain
`redirect-all-outgoing`-- as the name implies, it will be sent to the proxy to
be processed, on its outbound port `4140`. Much like in the inbound case, the
routing happens at the `nat` level, the packet's header will be re-written to
target the outbound port. The proxy will process the packet and then forward it
to its destination. The new packet will take the same path through the `OUTPUT`
chain, however, it will stop at the first rule, since it was produced by the
proxy.

The substantiated explanation applies to a packet whose destination is another
service, outside of the pod. In practice, an application can also send traffic
locally. As such, there are two other possible scenarios that we will explore:
_when a service talks to itself_ (by sending traffic over localhost or by using
its own endpoint address), and when _a service talks to itself through a
`clusterIP` target_. Both scenarios are somehow related, but the path a packet
takes differs.

**A service may send requests to itself**. It can also target another container
in the pod. This scenario would typically apply when:

* The destination is the pod (or endpoint) IP address.
* The destination is a port bound on localhost (regardless of which container
it belongs to).

![Outbound iptables chain traversal](/docs/images/iptables/iptables-fig2-3.png "Outbound iptables chain traversal")

When the application targets itself through its pod's IP (or loopback address),
the packets will traverse the two output chains. The first rule will be
skipped, since the owner is the application, and not the proxy. Once the second
rule is matched, the packets will return to the first output chain, from here,
they'll be sent directly to the service.

{{< note >}}
Usually, packets traverse another chain on the outbound side called
`POSTROUTING`. This chain is traversed after the `OUTPUT` chain, but to keep
the explanation simple, it has not been mentioned. Likewise, outbound packets that
are sent over the loopback interface become inbound packets, since they need to
be processed again. The kernel takes shortcuts in this case and bypasses the
`PREROUTING` chain that inbound packets from the outside world traverse when
they first arrive. For this reason, we do not need any special rules on the
inbound side to account for outbound packets that are sent locally.
{{< /note >}}

**A service may send requests to itself using its clusterIP**. In such cases,
it is not guaranteed that the destination will be local. The packet follows an
unusual path, as depicted in the diagram below.

![Outbound iptables chain traversal](/docs/images/iptables/iptables-fig2-4.png "Outbound iptables chain traversal")

When the packet first traverses the output chains, it will follow the same path
an outbound packet would normally take. In such a scenario, the packet's
destination will be an address that is not considered to be local by the
kernel-- it is, after all, a virtual IP. The proxy will process the packet, at
a connection level, connections to a `clusterIP` will be load balanced between
endpoints. Chances are that the endpoint selected will be the pod itself,
packets will therefore never leave the pod; the destination will be resolved to
the podIP. The packets produced by the proxy will traverse the output chain and
stop at the first rule, then they will be forwarded to the service. This
constitutes an edge case because at this point, the packet has been processed
by the proxy, unlike the scenario previously discussed where it skips it
altogether. For this reason, at a connection level, the proxy will _not_ mTLS
or opportunistically upgrade the connection to HTTP/2 when the endpoint is
local to the pod. In practice, this is treated as if the destination was
loopback, with the exception that the packet is forwarded through the proxy,
instead of being forwarded from the service directly to itself.

## Rules table

For reference, you can find the actual commands used to create the rules below.
Alternatively, if you want to inspect the iptables rules created for a pod, you
can retrieve them through the following command:

```bash
$ kubectl -n <namesppace> logs <pod-name> linkerd-init
# where <pod-name> is the name of the pod
# you want to see the iptables rules for
```
<!-- markdownlint-disable MD013 -->
### Inbound

{{< keyval >}}
| # | name | iptables rule | description|
|---|------|---------------|------------|
| 1 | redirect-common-chain | `iptables -t nat -N PROXY_INIT_REDIRECT`| creates a new `iptables` chain to add inbound redirect rules to; the chain is attached to the `nat` table |
| 2 | ignore-port | `iptables -t nat -A PROXY_INIT_REDIRECT -p tcp --match multiport --dports <ports> -j RETURN` | configures `iptables` to ignore the redirect chain for packets whose dst ports are included in the `--skip-inbound-ports` config option |
| 3 | proxy-init-redirect-all | `iptables -t nat -A PROXY_INIT_REDIRECT -p tcp -j REDIRECT --to-port 4143` | configures `iptables` to redirect all incoming TCP packets to port `4143`, the proxy's inbound port |
| 4 | install-proxy-init-prerouting | `iptables -t nat -A PREROUTING -j PROXY_INIT_REDIRECT` | the last inbound rule configures the `PREROUTING` chain (first chain a packet traverses inbound) to send packets to the redirect chain for processing |
{{< /keyval >}}

### Outbound

{{< keyval >}}
| # | name | iptables rule | description |
|---|------|---------------|-------------|
| 1 | redirect-common-chain | `iptables -t nat -N PROXY_INIT_OUTPUT`| creates a new `iptables` chain to add outbound redirect rules to, also attached to the `nat` table |
| 2 | ignore-proxy-uid | `iptables -t nat -A PROXY_INIT_OUTPUT -m owner --uid-owner 2102 -j RETURN` | when a packet is owned by the proxy (`--uid-owner 2102`), skip processing and return to the previous (`OUTPUT`) chain |
| 3 | ignore-loopback | `iptables -t nat -A PROXY_INIT_OUTPUT -o lo -j RETURN` | when a packet is sent over the loopback interface (`lo`), skip processing and return to the previous chain |
| 4 | ignore-port | `iptables -t nat -A PROXY_INIT_OUTPUT -p tcp --match multiport --dports <ports> -j RETURN` | configures `iptables` to ignore the redirect output chain for packets whose dst ports are included in the `--skip-outbound-ports` config option |
| 5 | redirect-all-outgoing | `iptables -t nat -A PROXY_INIT_OUTPUT -p tcp -j REDIRECT --to-port 4140`|  configures `iptables` to redirect all outgoing TCP packets to port `4140`, the proxy's outbound port |
| 6 | install-proxy-init-output | `iptables -t nat -A OUTPUT -j PROXY_INIT_OUTPUT` | the last outbound rule configures the `OUTPUT` chain (second before last chain a packet traverses outbound) to send packets to the redirect output chain for processing |
{{< /keyval >}}
<!-- markdownlint-enable MD013 -->
