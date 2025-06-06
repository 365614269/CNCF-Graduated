---
title: Quickstart
weight: 900
description: Get etcd up and running in less than 5 minutes!
---

Follow these instructions to locally install, run, and test a single-member
cluster of etcd:

 1. Install etcd from pre-built binaries or from source. For details, see
    [Install][].

    {{% alert color="warning" %}}**Important**: Ensure that you perform the last
    step of the installation instructions to verify that `etcd` is in your path.
    {{% /alert %}}

 2. Launch `etcd`:

    ```console
    $ etcd
    {"level":"info","ts":"2021-09-17T09:19:32.783-0400","caller":"etcdmain/etcd.go:72","msg":... }
    ⋮
    ```

    {{% alert color="info" %}}**Note**: The output produced by `etcd` are
    [logs](../op-guide/configuration/#logging) &mdash; info-level logs can
    be ignored. {{% /alert %}}

 3. From **another terminal**, use `etcdctl` to set a key:

    ```console
    $ etcdctl put greeting "Hello, etcd"
    OK
    ```

 4. From the same terminal, retrieve the key:

    ```console
    $ etcdctl get greeting
    greeting
    Hello, etcd
    ```

## What's next?

Learn about more ways to configure and use etcd from the following pages:

- If you are a developer:
  - Explore the gRPC [API][].
  - Find [language bindings and tools][integrations].

- If you are an operator or admin:
  - Set up a [multi-machine cluster][clustering].
  - Learn how to [configure][] etcd.
  - Use TLS to [secure an etcd cluster][security].
  - [Tune etcd][tuning].

[api]: /docs/{{< param version >}}/learning/api
[clustering]: /docs/{{< param version >}}/op-guide/clustering
[configure]: /docs/{{< param version >}}/op-guide/configuration
[integrations]: /docs/{{< param version >}}/integrations
[security]: /docs/{{< param version >}}/op-guide/security
[tuning]: /docs/{{< param version >}}/tuning
[Install]: ../install/
