---
title: Zipkin
description: 如何与 Zipkin 进行集成。
weight: 31
keywords: [integration,zipkin,tracing]
owner: istio/wg-environments-maintainers
test: n/a
---

[Zipkin](https://zipkin.io/) 是一个分布式追踪系统，可用来协助收集在定位服务架构延迟问题时所需的计时数据，
功能包括此数据的收集和查找。

## 安装 {#installation}

### 方法1：快速开始 {#quick-start}

Istio 提供了一个基本的安装示例来快速启动和运行 Zipkin：

{{< text bash >}}
$ kubectl apply -f {{< github_file >}}/samples/addons/extras/zipkin.yaml
{{< /text >}}

通过 `kubectl apply -f` 将 Zipkin 部署到集群中。此示例仅用于演示，
没有针对其性能或安全性进行调优。

### 方法2：定制化安装 {#customizable-install}

参阅 [Zipkin 文档](https://zipkin.io/)开始安装。Zipkin 集成 Istio 使用时无需进行特殊修改。

## 使用 {#usage}

有关使用 Zipkin 的信息，请参阅
[Zipkin](/zh/docs/tasks/observability/distributed-tracing/zipkin/)。
