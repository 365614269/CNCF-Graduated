---
title: 抽象
status: Completed
category: 屬性
tags: ["基本原理", "", ""]
---

在計算機的背景中，抽象是一種表示方式，它將細節隱藏起來，
讓[服務](/zh-tw/service/)的使用者（包括電腦程式和人類）能夠更容易理解系統並使其更通用。
一個很好的例子是您的筆記型電腦的作業系統（OS）。
它抽象化您的計算機運作的所有細節。
您不需要了解 CPU、記憶體以及程式如何運作，
您只需操作作業系統，作業系統會處理這些細節。
所有細節都被隱藏在作業系統的「幕後」或抽象之後。

系統通常具有多個抽象層。
這顯著簡化了開發過程。
在程式設計時，開發人員建立與特定抽象層兼容的元件，
而不需要擔心底層的具體細節及差異。
不論底層是什麼，只要與抽象層兼容，元件就能在系統中運作。
