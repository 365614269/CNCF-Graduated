---
title: 虛擬機器
status: Completed
category: 技術
tags: ["基本原理", "基礎設施", ""]
---

## 是什麼 {#what-it-is}

虛擬機器（VM）是一台不受限於特定硬體的計算機及其作業系統。
透過[虛擬化](/zh-tw/virtualization/)技術，
我們可以將一台實體的計算機分割成多台虛擬的計算機。
這種分割可以讓組織和基礎設施提供者能輕鬆地建立和刪除虛擬機器，
而不會影響底層的硬體。

## 解決的問題 {#problem-it-addresses}

當一台[裸機](/zh-tw/bare-metal-machine/)綁定到單一作業系統時，
該計算機的資源使用情況通常在某種程度上會受到限制。
另外，當一個作業系統綁定到單一實體計算機時，
其可用性也會與該硬體的直接相關。
如果實體計算機因為維護或硬體故障而離線，
該作業系統也會隨之離線。

## 如何幫助我們 {#how-it-helps}

透過解除作業系統與單一實體計算機之間的直接綁定，
我們便可以解決裸機的幾個問題：
佈建時間、硬體使用率和復原力。

在佈建新的虛擬機器時，
我們不需要為此購買、安裝或設定新的硬體，
所以佈建一台新計算機的時間便可以大幅地縮短。
由於我們可以將多台虛擬機器放置在一台實體計算機上，
這讓我們能夠更有效率地使用既有的硬體資源。
虛擬機器不受特定實體計算機的限制，
因此也比實體計算機更具有復原力。
當一台實體計算機需要離線時，
運行於其上的虛擬機器可以無需或只需極少的停機時間，
便可以轉移到另一台實體計算機上。
