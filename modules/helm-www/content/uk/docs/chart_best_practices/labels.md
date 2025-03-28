---
title: "Мітки та Анотації"
description: "Охоплює найкращі практики використання міток та анотацій у вашому Chart."
weight: 5
---

Цей розділ посібника з найкращих практик обговорює найкращі практики для використання міток та анотацій у вашому чарті.

## Це мітка чи анотація? {#is-it-a-label-or-an-annotation}

Елемент метаданих слід вважати міткою за таких умов:

- Він використовується Kubernetes для ідентифікації цього ресурсу
- Він корисний для експонування операторам з метою запиту системи.

Наприклад, рекомендується використовувати `helm.sh/chart: NAME-VERSION` як мітку, щоб оператори могли зручно знаходити всі екземпляри конкретного чарту.

Якщо елемент метаданих не використовується для запитів, його слід встановити як анотацію.

Helm hooks завжди є анотаціями.

## Стандартні Мітки {#standard-labels}

Наступна таблиця визначає загальні мітки, які використовуються в Helm чарті. Helm сам по собі ніколи не вимагає наявності конкретної мітки. Мітки, які позначені як REC, є рекомендованими та _повинні_ бути розміщені в чарті для глобальної узгодженості. Ті, що позначені як OPT, є необовʼязковими. Це ідіоматичні або часто використовувані мітки, але не є критично важливими для операційних цілей.

| Назва | Статус | Опис |
|-------|--------|------|
| `app.kubernetes.io/name` | REC | Це повинно бути імʼя застосунку, яке відображає весь застосунок. Зазвичай використовується `{{ template "name" . }}`. Це використовується багатьма маніфестами Kubernetes і не є специфічним для Helm. |
| `helm.sh/chart` | REC | Це повинно бути імʼя чарту та версія: `{{ .Chart.Name }}-{{ .Chart.Version \| replace "+" "_" }}`. |
| `app.kubernetes.io/managed-by` | REC | Це завжди повинно бути встановлено як `{{ .Release.Service }}`. Це для знаходження всього, що управляється Helm. |
| `app.kubernetes.io/instance` | REC | Це повинно бути `{{ .Release.Name }}`. Це допомагає відрізняти різні екземпляри одного й того ж застосунку. |
| `app.kubernetes.io/version` | OPT | Версія застосунку і може бути встановлена як `{{ .Chart.AppVersion }}`. |
| `app.kubernetes.io/component` | OPT | Це загальна мітка для позначення різних ролей, які елементи можуть відігравати в застосунку. Наприклад, `app.kubernetes.io/component: frontend`. |
| `app.kubernetes.io/part-of` | OPT | Коли кілька чартів або програмних частин використовуються разом для створення одного застосунку. Наприклад, програмне забезпечення та база даних для створення вебсайту. Це можна встановити на рівні основного застосунку, що підтримується. |

Ви можете знайти більше інформації про мітки Kubernetes, що починаються з `app.kubernetes.io`, в [документації Kubernetes](https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/).
