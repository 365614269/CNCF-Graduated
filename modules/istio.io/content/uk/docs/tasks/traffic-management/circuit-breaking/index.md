---
title: Розмикання ланцюга
description: У цьому завданні показано, як налаштувати розмикання ланцюга (запобіжника) для зʼєднань, запитів і виявлення відхилень.
weight: 50
keywords: [traffic-management,circuit-breaking]
owner: istio/wg-networking-maintainers
test: yes
---

У цьому завданні показано, як налаштувати розмикання ланцюга (запобіжника) для зʼєднань, запитів і виявлення відхилень.

Розмикання ланцюга є важливим шаблоном для створення відмовостійких мікросервісних застосунків. Розмикання ланцюгів дозволяє створювати застосунки, які обмежують вплив збоїв, стрибків затримок та інших небажаних ефектів, пов'язаних з особливостями мережі.

У цьому завданні ви налаштуєте правила розмикання ланцюга, а потім протестуєте конфігурацію, навмисно «вимкнувши» запобіжник.

## Перш ніж почати {#before-you-begin}

* Налаштуйте Istio, дотримуючись інструкцій у [керівництві з встановлення](/docs/setup/).

{{< boilerplate start-httpbin-service >}}

Застосунок `httpbin` виконує роль бекенд-сервісу для цього завдання.

## Налаштування запобіжника {#configuring-the-circuit-breaker}

1.  Створіть [правило призначення](/docs/reference/config/networking/destination-rule/), щоб застосувати налаштування розмикання ланцюга при виклику сервісу `httpbin`:

    {{< warning >}}
    Якщо ви встановили/налаштували Istio з увімкненою взаємною TLS автентифікацією, вам потрібно додати політику трафіку TLS `mode: ISTIO_MUTUAL` до `DestinationRule` перед його застосуванням. В іншому випадку запити будуть генерувати помилки 503, як описано [тут](/docs/ops/common-problems/network-issues/#503-errors-after-setting-destination-rule).
    {{< /warning >}}

    {{< text bash >}}
    $ kubectl apply -f - <<EOF
    apiVersion: networking.istio.io/v1
    kind: DestinationRule
    metadata:
      name: httpbin
    spec:
      host: httpbin
      trafficPolicy:
        connectionPool:
          tcp:
            maxConnections: 1
          http:
            http1MaxPendingRequests: 1
            maxRequestsPerConnection: 1
        outlierDetection:
          consecutive5xxErrors: 1
          interval: 1s
          baseEjectionTime: 3m
          maxEjectionPercent: 100
    EOF
    {{< /text >}}

2.  Переконайтесь, що правило було створено:

    {{< text bash yaml >}}
    $ kubectl get destinationrule httpbin -o yaml
    apiVersion: networking.istio.io/v1
    kind: DestinationRule
    ...
    spec:
      host: httpbin
      trafficPolicy:
        connectionPool:
          http:
            http1MaxPendingRequests: 1
            maxRequestsPerConnection: 1
          tcp:
            maxConnections: 1
        outlierDetection:
          baseEjectionTime: 3m
          consecutive5xxErrors: 1
          interval: 1s
          maxEjectionPercent: 100
    {{< /text >}}

## Додавання клієнта {#adding-a-client}

Створіть клієнта для надсилання трафіку до сервісу `httpbin`. Клієнт — це простий інструмент для навантажувального тестування з назвою [fortio](https://github.com/istio/fortio). Fortio дозволяє вам контролювати кількість зʼєднань, одночасність і затримки для вихідних HTTP запитів. Ви будете використовувати цей клієнт, щоб "вимкнути" політики розриву ланцюга, які ви встановили в `DestinationRule`.

1. Виконайте інʼєкцію клієнта з проксі-сервером Istio, щоб мережеві взаємодії керувалися Istio.

    Якщо ви увімкнули [автоматичну інʼєкцію sidecar](/docs/setup/additional-setup/sidecar-injection/#automatic-sidecar-injection), розгорніть сервіс `fortio`:

    {{< text bash >}}
    $ kubectl apply -f @samples/httpbin/sample-client/fortio-deploy.yaml@
    {{< /text >}}

    В іншому випадку вам потрібно вручну виконати інʼєкцію sidecar перед розгортанням застосунку `fortio`:

    {{< text bash >}}
    $ kubectl apply -f <(istioctl kube-inject -f @samples/httpbin/sample-client/fortio-deploy.yaml@)
    {{< /text >}}

2. Увійдіть до podʼа клієнта і використовуйте інструмент fortio для виклику `httpbin`. Скористайтесь `curl`, щоб вказати, що ви хочете зробити лише один виклик:

    {{< text bash >}}
    $ export FORTIO_POD=$(kubectl get pods -l app=fortio -o 'jsonpath={.items[0].metadata.name}')
    $ kubectl exec "$FORTIO_POD" -c fortio -- /usr/bin/fortio curl -quiet http://httpbin:8000/get
    HTTP/1.1 200 OK
    server: envoy
    date: Tue, 25 Feb 2020 20:25:52 GMT
    content-type: application/json
    content-length: 586
    access-control-allow-origin: *
    access-control-allow-credentials: true
    x-envoy-upstream-service-time: 36

    {
      "args": {},
      "headers": {
        "Content-Length": "0",
        "Host": "httpbin:8000",
        "User-Agent": "fortio.org/fortio-1.3.1",
        "X-B3-Parentspanid": "8fc453fb1dec2c22",
        "X-B3-Sampled": "1",
        "X-B3-Spanid": "071d7f06bc94943c",
        "X-B3-Traceid": "86a929a0e76cda378fc453fb1dec2c22",
        "X-Forwarded-Client-Cert": "By=spiffe://cluster.local/ns/default/sa/httpbin;Hash=68bbaedefe01ef4cb99e17358ff63e92d04a4ce831a35ab9a31d3c8e06adb038;Subject=\"\";URI=spiffe://cluster.local/ns/default/sa/default"
      },
      "origin": "127.0.0.1",
      "url": "http://httpbin:8000/get"
    }
    {{< /text >}}

Ви можете бачити, що запит пройшов успішно! Тепер настав час щось зламати.

## Використання запобіжника {#tripping-the-circuit-breaker}

У налаштуваннях `DestinationRule` ви вказали `maxConnections: 1` і `http1MaxPendingRequests: 1`. Ці правила вказують на те, що якщо у вас більше ніж одне одночасне зʼєднання та запит, ви повинні побачити деякі відмови, коли `istio-proxy` відкриє ланцюг для подальших запитів та зʼєднань.

1. Викличте сервіс з двома одночасними зʼєднаннями (`-c 2`) і надішліть 20 запитів (`-n 20`):

    {{< text bash >}}
    $ kubectl exec "$FORTIO_POD" -c fortio -- /usr/bin/fortio load -c 2 -qps 0 -n 20 -loglevel Warning http://httpbin:8000/get
    20:33:46 I logger.go:97> Log level is now 3 Warning (was 2 Info)
    Fortio 1.3.1 running at 0 queries per second, 6->6 procs, for 20 calls: http://httpbin:8000/get
    Starting at max qps with 2 thread(s) [gomax 6] for exactly 20 calls (10 per thread + 0)
    20:33:46 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:33:47 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:33:47 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    Ended after 59.8524ms : 20 calls. qps=334.16
    Aggregated Function Time : count 20 avg 0.0056869 +/- 0.003869 min 0.000499 max 0.0144329 sum 0.113738
    # range, mid point, percentile, count
    >= 0.000499 <= 0.001 , 0.0007495 , 10.00, 2
    > 0.001 <= 0.002 , 0.0015 , 15.00, 1
    > 0.003 <= 0.004 , 0.0035 , 45.00, 6
    > 0.004 <= 0.005 , 0.0045 , 55.00, 2
    > 0.005 <= 0.006 , 0.0055 , 60.00, 1
    > 0.006 <= 0.007 , 0.0065 , 70.00, 2
    > 0.007 <= 0.008 , 0.0075 , 80.00, 2
    > 0.008 <= 0.009 , 0.0085 , 85.00, 1
    > 0.011 <= 0.012 , 0.0115 , 90.00, 1
    > 0.012 <= 0.014 , 0.013 , 95.00, 1
    > 0.014 <= 0.0144329 , 0.0142165 , 100.00, 1
    # target 50% 0.0045
    # target 75% 0.0075
    # target 90% 0.012
    # target 99% 0.0143463
    # target 99.9% 0.0144242
    Sockets used: 4 (for perfect keepalive, would be 2)
    Code 200 : 17 (85.0 %)
    Code 503 : 3 (15.0 %)
    Response Header Sizes : count 20 avg 195.65 +/- 82.19 min 0 max 231 sum 3913
    Response Body/Total Sizes : count 20 avg 729.9 +/- 205.4 min 241 max 817 sum 14598
    All done 20 calls (plus 0 warmup) 5.687 ms avg, 334.2 qps
    {{< /text >}}

    Цікаво, що майже всі запити були виконані! Проксі-сервер `istio-proxy` дійсно дозволяє деяку свободу дій.

    {{< text plain >}}
    Code 200 : 17 (85.0 %)
    Code 503 : 3 (15.0 %)
    {{< /text >}}

1. Доведіть кількість одночасних зʼєднань до 3:

    {{< text bash >}}
    $ kubectl exec "$FORTIO_POD" -c fortio -- /usr/bin/fortio load -c 3 -qps 0 -n 30 -loglevel Warning http://httpbin:8000/get
    20:32:30 I logger.go:97> Log level is now 3 Warning (was 2 Info)
    Fortio 1.3.1 running at 0 queries per second, 6->6 procs, for 30 calls: http://httpbin:8000/get
    Starting at max qps with 3 thread(s) [gomax 6] for exactly 30 calls (10 per thread + 0)
    20:32:30 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:32:30 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:32:30 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:32:30 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:32:30 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:32:30 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:32:30 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:32:30 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:32:30 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:32:30 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:32:30 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:32:30 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:32:30 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:32:30 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:32:30 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:32:30 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:32:30 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:32:30 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    20:32:30 W http_client.go:679> Parsed non ok code 503 (HTTP/1.1 503)
    Ended after 51.9946ms : 30 calls. qps=576.98
    Aggregated Function Time : count 30 avg 0.0040001633 +/- 0.003447 min 0.0004298 max 0.015943 sum 0.1200049
    # range, mid point, percentile, count
    >= 0.0004298 <= 0.001 , 0.0007149 , 16.67, 5
    > 0.001 <= 0.002 , 0.0015 , 36.67, 6
    > 0.002 <= 0.003 , 0.0025 , 50.00, 4
    > 0.003 <= 0.004 , 0.0035 , 60.00, 3
    > 0.004 <= 0.005 , 0.0045 , 66.67, 2
    > 0.005 <= 0.006 , 0.0055 , 76.67, 3
    > 0.006 <= 0.007 , 0.0065 , 83.33, 2
    > 0.007 <= 0.008 , 0.0075 , 86.67, 1
    > 0.008 <= 0.009 , 0.0085 , 90.00, 1
    > 0.009 <= 0.01 , 0.0095 , 96.67, 2
    > 0.014 <= 0.015943 , 0.0149715 , 100.00, 1
    # target 50% 0.003
    # target 75% 0.00583333
    # target 90% 0.009
    # target 99% 0.0153601
    # target 99.9% 0.0158847
    Sockets used: 20 (for perfect keepalive, would be 3)
    Code 200 : 11 (36.7 %)
    Code 503 : 19 (63.3 %)
    Response Header Sizes : count 30 avg 84.366667 +/- 110.9 min 0 max 231 sum 2531
    Response Body/Total Sizes : count 30 avg 451.86667 +/- 277.1 min 241 max 817 sum 13556
    All done 30 calls (plus 0 warmup) 4.000 ms avg, 577.0 qps
    {{< /text >}}

    Тепер ви починаєте бачити очікувану поведінку запобіжника. Лише 36.7% запитів були успішними, а решта була перехоплена запобіжником:

    {{< text plain >}}
    Code 200 : 11 (36.7 %)
    Code 503 : 19 (63.3 %)
    {{< /text >}}

2.  Щоб дізнатися більше, перегляньте статистику `istio-proxy`:

    {{< text bash >}}
    $ kubectl exec "$FORTIO_POD" -c istio-proxy -- pilot-agent request GET stats | grep httpbin | grep pending
    cluster.outbound|8000||httpbin.default.svc.cluster.local;.circuit_breakers.default.remaining_pending: 1
    cluster.outbound|8000||httpbin.default.svc.cluster.local;.circuit_breakers.default.rq_pending_open: 0
    cluster.outbound|8000||httpbin.default.svc.cluster.local;.circuit_breakers.high.rq_pending_open: 0
    cluster.outbound|8000||httpbin.default.svc.cluster.local;.upstream_rq_pending_active: 0
    cluster.outbound|8000||httpbin.default.svc.cluster.local;.upstream_rq_pending_failure_eject: 0
    cluster.outbound|8000||httpbin.default.svc.cluster.local;.upstream_rq_pending_overflow: 21
    cluster.outbound|8000||httpbin.default.svc.cluster.local;.upstream_rq_pending_total: 29
    {{< /text >}}

    Ви можете побачити `21` для значення `upstream_rq_pending_overflow`, що означає, що поки що `21` викликів було позначено для обробки запобіжником.

## Очищення {#cleaning-up}

1.  Видаліть правило:

    {{< text bash >}}
    $ kubectl delete destinationrule httpbin
    {{< /text >}}

2.  Припиніть роботу сервісу [httpbin]({{< github_tree >}}/samples/httpbin) та клієнта:

    {{< text bash >}}
    $ kubectl delete -f @samples/httpbin/sample-client/fortio-deploy.yaml@
    $ kubectl delete -f @samples/httpbin/httpbin.yaml@
    {{< /text >}}
