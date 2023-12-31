---
title: 抽象
status: Completed
category: 属性
tags: ["架构", "", ""]
---

在计算的上下文中，抽象是一种对 [服务](/service/) 的消费者（消费者可以是计算机程序或人类）隐藏其细节的表示法，使系统更通用也更容易理解。
你的笔记本电脑的操作系统就是一个很好的例子。它把计算机工作的所有细节都抽象出来了。
你不需要知道任何关于CPU、内存以及程序如何被运行，你只需操作操作系统，操作系统会处理这些细节。
所有这些细节都隐藏在操作系统的“幕布”或抽象概念后面。

系统通常有多个抽象层。这大大简化了开发工作。在编程时，开发人员构建与特定抽象层兼容的组件，而不必担心可能非常异构的所有底层细节。
如果组件能与抽象层一起工作，它就能与系统一起工作 —— 无论底层是什么样的。
