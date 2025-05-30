---
title: Supported Kubernetes Versions
description: Reference documentation for which Linkerd version supports which Kubernetes
  version
---

Linkerd supports all versions of Kubernetes that were supported at the time
that a given Linkerd version ships. For example, at the time that Linkerd 2.14
shipped, Kubernetes versions 1.26, 1.27, and 1.28 were supported, so Linkerd
2.14 supports all of those Kubernetes versions. (In many cases, as you'll see
below, Linkerd versions will also support older Kubernetes versions.)

Obviously, Linkerd 2.14 has no knowledge of what changes will come _after_
Kubernetes 1.28. In some cases, later versions of Kubernetes end up making
changes that cause older versions of Linkerd to not work: we will update the
chart below as these situations arise.

{{< keyval >}}
| Linkerd Version | Minimum Kubernetes Version | Maximum Kubernetes Version |
|-----------------|----------------------------|----------------------------|
| `2.10`          | `1.16`                     | `1.23`                     |
| `2.11`          | `1.17`                     | `1.23`                     |
| `2.12`          | `1.21`                     | `1.24`                     |
| `2.13`          | `1.21`                     | `1.28`                     |
| `2.14`          | `1.21`                     | `1.28`                     |
| `2.15`          | `1.22`                     | `1.29`                     |
| `2.16`          | `1.22`                     | `1.29`                     |
| `2.17`          | `1.22`                     | `1.31`                     |
| `2.18`          | `1.22`                     | `1.32`                     |
{{< /keyval >}}

Note that Linkerd will almost never change the supported Kubernetes version in
a minor release, which is why the table above only lists major versions. One
known exception: Linkerd 2.11.0 supported Kubernetes 1.16, but 2.11.1 and
later required Kubernetes 1.17 as shown in the table above.

## Edge Releases

{{< keyval >}}
| Linkerd Version | Minimum Kubernetes Version | Maximum Kubernetes Version |
|-----------------|----------------------------|----------------------------|
| `edge-22.10.1` - `edge-23.12.1` | `1.21`     | `1.29`                     |
| `edge-23.12.2` - `edge-24.11.4` | `1.22`     | `1.29`                     |
| `edge-24.11.5` - `edge-25.1.2`  | `1.22`     | `1.31`                     |
| `edge-25.2.1` and newer         | `1.22`     | `1.32`                     |
{{< /keyval >}}
