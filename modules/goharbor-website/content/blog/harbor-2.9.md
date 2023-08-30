---
title: "Harbor v2.9 release"
author:
  name: "Orlin Vasilev"
date: 2023-08-29T06:00:00+01:00
draft: false
showPageInfo: true
---

# Harbor 2.9

##### Introducing [Harbor 2.9][r29]
![banner](../img/harbor-2.9.png)
This release comes with several new features and improvements, including Harbor Security Hub and OCI v1.1.0-rc2, which provides enhanced support for multi-architecture images and image manifests. Another major addition is improved Garbage Collector. [The 2.9 release][r29] also includes an Customizable Message Banner(check green banner above) , providing more comprehensive and detailed information about upcoming maintenance and admin activities. Also we have introduces some significant changes as the removal of Notary.

## New features

### Harbor Security Hub
Admin users can now access valuable security insights, which including the number of scanned and unscanned artifacts, identification of dangerous artifacts and CVEs and advanced search capabilities for vulnerabilities using multiple combined conditions.

* Add Security Hub UI by @AllForNothing in [18942](https://github.com/goharbor/harbor/pull/18942)
* Update table scan_report and extract cvss_v3_score from vendor attribute by @stonezdj in [18854](https://github.com/goharbor/harbor/pull/18854)
* Add vulnerability search API by @stonezdj in [18924](https://github.com/goharbor/harbor/pull/18924)
* Add security hub summary API by @stonezdj in [18872](https://github.com/goharbor/harbor/pull/18872)
* Create index in vulnerability_record table by @stonezdj in [18949](https://github.com/goharbor/harbor/pull/18949)

[Harbor Security Hub Documentation](https://goharbor.io/docs/2.9.0/administration/security-hub/)

### GC Enhancements
Improved visibility with detailed GC execution history and enable parallel deletion for faster GC triggers.

* Add worker parameter for GC by @AllForNothing in [18882](https://github.com/goharbor/harbor/pull/18882)
* add more details in gc history by @wy65701436 in [18779](https://github.com/goharbor/harbor/pull/18779)
* add multiple deletion of GC by @wy65701436 in [18855](https://github.com/goharbor/harbor/pull/18855)

### [OCI Distribution Spec v1.1.0-rc2 Support](https://github.com/opencontainers/distribution-spec/releases/tag/v1.1.0-rc2)

Harbor now supports OCI Distribution Spec v1.1.0-rc2 and added support for **Notation** signature and **Nydus** conversion as referrers. This is a significant upgrade that allows users to store and distribute OCI and Docker images. With this feature, users can manage their images more effectively and efficiently.

* support OCI-Subject header by @wy65701436 in [18885](https://github.com/goharbor/harbor/pull/18885)
* add notation support by @wy65701436 in [18909](https://github.com/goharbor/harbor/pull/18909)
* enable notary v2 policy checker by @wy65701436 in [18927](https://github.com/goharbor/harbor/pull/18927)
* Add Notation UI for deployment security by @AllForNothing in [19952](https://github.com/goharbor/harbor/pull/18952)
* support nydus as a accessory by @wy65701436 in [18953](https://github.com/goharbor/harbor/pull/18953)




## Removal of Notary
Starting in Harbor 2.6 was decided to deprecate Notary you can check this [discussion](https://github.com/goharbor/harbor/discussions/16612). Workflows and pipelines needs to be migrated to [cosign as of Harbor 2.5](https://goharbor.io/blog/cosign-2.5.0/), with version v2.9.0, Harbor no longer includes Notary in either the user interface or the backend.

https://github.com/goharbor/harbor/pull/18620  
https://github.com/goharbor/harbor/pull/18666  
https://github.com/goharbor/harbor/pull/18668  

### Other Enhancements and Updates
* Fix message prompt under the header by @AllForNothing in https://github.com/goharbor/harbor/pull/18613
* fix: improve the performance of list artifacts by @chlins in https://github.com/goharbor/harbor/pull/18610
* Improve repo_read_only header on the UI by @AllForNothing in https://github.com/goharbor/harbor/pull/18729
* Add a text to explain the time window for GC by @AllForNothing in https://github.com/goharbor/harbor/pull/18735
* Add a tooltip for slack notification by @AllForNothing in https://github.com/goharbor/harbor/pull/18787
*【UT】add unit test for collector system info by @lengrongfu in https://github.com/goharbor/harbor/pull/18717
* Add Details column for gc history by @AllForNothing in https://github.com/goharbor/harbor/pull/18797
* Add Podman push command to the UI by @AllForNothing in https://github.com/goharbor/harbor/pull/18810
* Add new client Podman to the pull command by @AllForNothing in https://github.com/goharbor/harbor/pull/18857


##### If you're Harbor user and want to share any feedback, we'd love to hear from you [here](https://github.com/goharbor/community/issues/115)!

## Contributions to v2.9
# **Huge shout-out**
to the maintainer team, our contributors & everyone in the
community who helped the project realize its level of adoption today !!
The project would not be where it is today without the community


We also want to thank the following users for their sustained
contributions to the project!

Thank you to all who contributed to 2.9 [full release notes here][r29]


### Special thank you to all new contributors:
* @pgillich made their first contribution in https://github.com/goharbor/harbor/pull/18181
* @Maxi-Mega made their first contribution in https://github.com/goharbor/harbor/pull/18151
* @yrs147 made their first contribution in https://github.com/goharbor/harbor/pull/18282
* @perjahn made their first contribution in https://github.com/goharbor/harbor/pull/18588
* @dyf991645 made their first contribution in https://github.com/goharbor/harbor/pull/18602
* @PeterDaveHello made their first contribution in https://github.com/goharbor/harbor/pull/18605
* @iAklis made their first contribution in https://github.com/goharbor/harbor/pull/18139
* @DavidSpek made their first contribution in https://github.com/goharbor/harbor/pull/18685
* @malmor made their first contribution in https://github.com/goharbor/harbor/pull/18659
* @mcsage made their first contribution in https://github.com/goharbor/harbor/pull/16990
* @lishaokai1995 made their first contribution in https://github.com/goharbor/harbor/pull/18746
* @orblazer made their first contribution in https://github.com/goharbor/harbor/pull/18188
* @cuishuang made their first contribution in https://github.com/goharbor/harbor/pull/18726
* @sll552 made their first contribution in https://github.com/goharbor/harbor/pull/18783
* @vndroid made their first contribution in https://github.com/goharbor/harbor/pull/18767

## Collaborate with the Harbor Community

Get updates on Twitter: [@project\_harbor](https://twitter.com/project_harbor)  
Chat with us on Slack: [#harbor](https://cloud-native.slack.com/messages/harbor)
and [#harbor-dev](https://cloud-native.slack.com/messages/harbor-dev)
on the[CNCF Slack](https://slack.cncf.io)  
Collaborate with us on [GitHub](https://github.com/goharbor/harbor)  
Attend the [community meetings](https://github.com/goharbor/community/wiki/Harbor-Community-Meetings)  

&nbsp;
&nbsp;

Orlin Vasilev  
Harbor Community/DevRel  
GitHub: [@OrlinVasilev](https://github.com/OrlinVasilev)  
Twitter: [@OrlinVasilev](https://twitter.com/OrlinVasilev)


[r29]: https://github.com/goharbor/harbor/releases/tag/v2.9.0
