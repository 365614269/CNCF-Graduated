---
title: WatchBookmark
content_type: feature_gate
_build:
  list: never
  render: false

stages:
  - stage: alpha 
    defaultValue: false
    fromVersion: "1.15"
    toVersion: "1.15"
  - stage: beta
    defaultValue: true
    fromVersion: "1.16"  
    toVersion: "1.16" 
  - stage: stable
    defaultValue: true
    fromVersion: "1.17"  
    toVersion: "1.32"

removed: true
---
Enable support for watch bookmark events.
