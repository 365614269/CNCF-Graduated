---
title: kube-scheduler
id: kube-scheduler
date: 2018-04-12
full_link: /docs/reference/generated/kube-scheduler/
short_description: >
  Компонент управляющего слоя, который отслеживает недавно созданные поды без назначенного для них узла и выбирает узел, на котором они должны работать.

aka:
tags:
- architecture
---
 Компонент управляющего слоя (control plane), который отслеживает недавно созданные поды без назначенного для них узла и выбирает узел, на котором они должны работать.

<!--more-->

При планировании учитывается множество факторов, включая индивидуальные
и общие требования к ресурсам, ограничения по железу/программному обеспечению/политикам,
конфигурацию принадлежности (affinity) и непринадлежности (anti-affinity)
узлов/подов, местонахождение данных, взаимодействие между рабочими
нагрузками и дедлайны.
