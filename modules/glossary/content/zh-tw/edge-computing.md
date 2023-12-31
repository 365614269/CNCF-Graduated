---
title: 邊緣運算
status: Completed
category: Technology
---

## 是什麼 {#what-it-is}

邊緣運算是個[分散式系統](/zh-tw/distributed-systems/)，它將一些儲存和運算資源從主要資料中心轉移到資料來源。
收集到的資料在本地端（例如：工廠、商店或整座城市）進行計算，而不是傳送到集中式資料中心進行處理和分析。
這些本地端處理單元或裝置代表系統的邊緣，而資料中心則代表系統的中心。
邊緣計算出的結果會被送回主要資料中心做進一步處理。
邊緣運算的例子包括手腕上的小配件或分析交通流量的電腦。

## 解決的問題 {#problem-it-addresses}

過去十年中，我們可以看到越來越多的邊緣裝置（例如：手機、智慧型手錶、感測器）。
在某些情況下，即時資料處理不僅是一個不錯的選擇，而且極其重要。
想想自動駕駛的汽車。
現在想像一下，汽車感測器的資料必須先傳送到資料中心進行處理，然後再送回汽車，這樣汽車才能作出適當的反應。
如此產生的網路延遲會是致命的。
雖然這是一個極端的例子，但大多數使用者都不願意使用無法即時反應的智慧設備。

## 如何幫助我們 {#how-it-helps}

如上所述，要使邊緣設備發揮作用，它們必須至少在本地端完成部分處理和分析工作，以便於對使用者提供接近即時的回饋。
要做到這點，就必須將資料中心的部分儲存和處理資源轉移到資料生產地：邊緣設備。
已處理和未處理的資料隨後發送到資料中心做進一步處理和儲存。
簡而言之，效率和速度是邊緣運算的主要驅動力。
