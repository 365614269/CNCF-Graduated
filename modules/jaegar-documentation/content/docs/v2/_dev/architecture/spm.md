---
title: Service Performance Monitoring (SPM)
aliases: [../spm]
hasparent: true
---

[![Service Performance Monitoring](/img/frontend-ui/spm.png)](/img/frontend-ui/spm.png)

Surfaced in Jaeger UI as the "Monitor" tab, the motivation for this feature is
to help identify interesting traces (e.g. high QPS, slow or erroneous requests)
without needing to know the service or operation names up-front.

It is essentially achieved through aggregating span data to produce RED
(Request, Error, Duration) metrics.

Potential use cases include:

- Post deployment sanity checks across the org, or on known dependent services
  in the request chain.
- Monitoring and root-causing when alerted of an issue.
- Better onboarding experience for new users of Jaeger UI.
- Long-term trend analysis of QPS, errors and latencies.
- Capacity planning.

## UI Feature Overview

The "Monitor" tab provides a service-level aggregation, as well as an operation-level
aggregation within the service, of Request rates, Error rates and Durations
(P95, P75 and P50), also known as RED metrics.

Within the operation-level aggregations, an "Impact" metric, computed as the
product of latency and request rate, is another signal that can be used to
rule-out operations that may naturally have a high latency profile such as daily
batch jobs, or conversely highlight operations that are lower in the latency
rankings but with a high RPS (request per second).

From these aggregations, Jaeger UI is able to pre-populate a Trace search with
the relevant service, operation and lookback period, narrowing down the search
space for these more interesting traces.

## Getting Started

{{< info >}}
This is for demonstration purposes only and does not reflect deployment best practices.
{{< /info >}}

A locally runnable setup is available in the [Jaeger repository][spm-demo] along
with instructions on how to run it.

The feature can be accessed from the "Monitor" tab along the top menu.

This demo includes [Microsim](https://github.com/yurishkuro/microsim); a microservices
simulator to generate trace data.

If generating traces manually is preferred, the [Sample App: HotROD](../../getting-started/#-hotrod-demo) can be started via docker. Be sure to include `--net monitor_backend` in the `docker run` command.

## Configuration

### Option 1: PromQL-compatible backend Configuration

An example configuration is available in the Jaeger repository: [config-spm.yaml](https://github.com/jaegertracing/jaeger/tree/main/cmd/jaeger/config-spm.yaml). The following steps are required to enable the SPM feature:

* Enable the [SpanMetrics Connector][spanmetrics-conn] in the pipeline:
```yaml
# Declare an exporter for metrics produced by the connector.
# For example, a Prometheus server may be configured to scrape
# the metrics from this endpoint.
exporters:
  prometheus:
    endpoint: "0.0.0.0:8889"

# Declare spanmetrics connector.
connectors:
  spanmetrics:
    # any connector configuration options
    ...

# Enable the spanmetrics connector to bridge
# the traces pipeline into the metrics pipeline.
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [jaeger_storage_exporter, spanmetrics]
    metrics/spanmetrics:
      receivers: [spanmetrics]
      exporters: [prometheus]
```
* Define a remote PromQL-compatible storage under `metric_backends:` in the `jaeger_storage` extension:
```yaml
extensions:
  jaeger_storage:
    backends:
      some_trace_storage:
        ...
    metric_backends:
      some_metrics_storage:
        prometheus:
          endpoint: http://prometheus:9090
```
* Reference this metrics store in the `jaeger_query` extension:
```yaml
extensions:
  jaeger_query:
    traces: some_trace_storage
    metrics_storage: some_metrics_storage
```
* Set the `monitor.menuEnabled=true` property in the [Jaeger UI configuration](../../deployment/frontend-ui/#monitor).

### Option 2: Elasticsearch/OpenSearch Backend Configuration

An example configuration is available in the Jaeger repository: [config-spm-elasticsearch.yaml](https://github.com/jaegertracing/jaeger/tree/main/cmd/jaeger/config-spm-elasticsearch.yaml).

Configuration is simpler as it does not require the SpanMetrics Connector or a separate metrics pipeline.

* Ensure your `jaeger_storage` extension is configured with an Elasticsearch or OpenSearch backend and metric_backends.

```yaml
extensions:
  jaeger_storage:
    backends:
      elasticsearch_trace_storage: &elasticsearch_config
        elasticsearch:
          ...
    metric_backends:
      elasticsearch_trace_storage: *elasticsearch_config # This points the metric backend to the same Elasticsearch trace storage configuration
```

```yaml
extensions:
  jaeger_storage:
    backends:
      opensearch_trace_storage: &opensearch_config
        opensearch:
          ...
    metric_backends:
      opensearch_trace_storage: *opensearch_config # This points the metric backend to the same OpenSearch trace storage configuration
```

* In the `jaeger_query` extension, set both `traces` and `metrics` storage to use the same Elasticsearch or OpenSearch backend defined in `jaeger_storage`.

```yaml
extensions:
  jaeger_query:
    storage:
      traces: elasticsearch_trace_storage # Must match the backend defined in jaeger_storage extension
      metrics: elasticsearch_trace_storage # Use the same storage for metrics
```

```yaml
extensions:
  jaeger_query:
    storage:
      traces: opensearch_trace_storage # Must match the backend defined in jaeger_storage extension
      metrics: opensearch_trace_storage # Use the same storage for metrics
```

* Set the `monitor.menuEnabled=true` property in the [Jaeger UI configuration](../../deployment/frontend-ui/#monitor).

## Architecture

There are two architectural approaches to generating RED metrics:

1.  **Pre-computing metrics**: using the [SpanMetrics Connector][spanmetrics-conn] to compute the metrics from spans and store them in a PromQL-compatible backend storage, which Jaeger queries.
2.  **Direct to storage**: Jaeger Query computes RED metrics at query time by directly querying the primary trace storage backend (Elasticsearch or OpenSearch). This simplifies the architecture by removing the need for a separate metrics pipeline and storage.

### Option 1: Pre-computing metrics

In addition to the standard Jaeger architecture, this approach requires:

- A [SpanMetrics Connector][spanmetrics-conn] is introduced in the pipeline that receives trace data (spans) and generates RED metrics.
- The generated metrics are exported to a Prometheus-compatible metrics store. In the provided example this is achieved by defining a `prometheus` exporter that opens an HTTP endpoint, and configuring a Prometheus server to scape the metrics from that endpoint. An alternative approach could be a push-style exporter that writes to a remote metrics store.
- An external Metrics Store that supports PromQL queries.
- A configuration in the `jaeger_query` extension to reference the external metrics store.

{{<mermaid align="center">}}
graph LR
    OTLP_EXPORTER[OTLP Exporter] --> TRACE_RECEIVER

    subgraph Application
        subgraph OpenTelemetry SDK
            OTLP_EXPORTER
        end
    end

    TRACE_RECEIVER[Trace Receiver] --> |spans| SPANMETRICS_CONN[SpanMetrics Connector]
    TRACE_RECEIVER --> |spans| TRACE_EXPORTER[Trace Exporter]
    TRACE_EXPORTER --> |spans| SPAN_STORE[(Trace Storage)]
    SPANMETRICS_CONN --> |metrics| PROMETHEUS_EXPORTER[Prometheus Exporter]
    PROMETHEUS_EXPORTER --> |metrics| METRICS_STORE[(Metrics Storage)]

    SPAN_STORE --> QUERY[Jaeger Query]
    METRICS_STORE --> QUERY
    QUERY --> UI[Jaeger UI]

    subgraph Jaeger all-in-one
        subgraph Pipeline
            TRACE_RECEIVER
            SPANMETRICS_CONN
            TRACE_EXPORTER
            PROMETHEUS_EXPORTER
            QUERY
            UI
        end
    end

    style Application fill:#DFDFDF,color:black

    style OTLP_EXPORTER fill:#404CA8,color:white
    style TRACE_RECEIVER fill:#404CA8,color:white
    style TRACE_EXPORTER fill:#404CA8,color:white
    style SPANMETRICS_CONN fill:#404CA8,color:white
    style PROMETHEUS_EXPORTER fill:#404CA8,color:white

    style UI fill:#9AEBFE,color:black
    style QUERY fill:#9AEBFE,color:black
{{< /mermaid >}}

### Option 2: Elasticsearch/OpenSearch Backend

This approach computes metrics directly from trace data stored in Elasticsearch or OpenSearch, eliminating the need for a separate metrics storage backend like Prometheus. The OpenTelemetry Collector is still used to receive traces and forward them to Jaeger, but the SpanMetrics Connector is not required.

{{<mermaid align="center">}}
graph LR
    OTLP_EXPORTER[OTLP Exporter] --> TRACE_RECEIVER

    subgraph Application
        subgraph OpenTelemetry SDK
            OTLP_EXPORTER
        end
    end

    TRACE_RECEIVER[Trace Receiver] --> |spans| TRACE_EXPORTER[Trace Exporter]
    TRACE_EXPORTER --> |spans| TRACE_STORAGE[(Elasticsearch / OpenSearch)]

    TRACE_STORAGE -->|traces & metrics queries| QUERY[Jaeger Query]
    QUERY --> UI[Jaeger UI]

    subgraph Jaeger all-in-one
        subgraph Pipeline
            TRACE_RECEIVER
            TRACE_EXPORTER
            QUERY
            UI
        end
    end

    style Application fill:#DFDFDF,color:black
    style TRACE_RECEIVER fill:#404CA8,color:white
    style TRACE_EXPORTER fill:#404CA8,color:white
    style UI fill:#9AEBFE,color:black
    style QUERY fill:#9AEBFE,color:black
{{< /mermaid >}}

### Metrics Storage

When using the PromQL-compatible backend architecture, any PromQL-compatible backend is supported by Jaeger Query. A list of these have been compiled by Julius Volz in: [https://promlabs.com/blog/2020/11/26/an-update-on-promql-compatibility-across-vendors](https://promlabs.com/blog/2020/11/26/an-update-on-promql-compatibility-across-vendors)

When using the direct querying architecture, **Elasticsearch** and **OpenSearch** are supported for both trace storage and metrics calculation.

### Derived Time Series

{{< info >}}
This section applies only to the **PromQL-compatible backend** architecture.
{{< /info >}}

It is worth understanding the additional metrics and time series that the
[SpanMetrics Connector][spanmetrics-conn] will generate in metrics storage to help
with capacity planning when deploying SPM.

Please refer to [Prometheus documentation][prom-metric-labels] covering the
concepts of metric names, types, labels and time series; terms that will be used
in the remainder of this section.

Two metric names will be created:
- `calls_total`
  - **Type**: counter
  - **Description**: counts the total number of spans, including error spans.
    Call counts are differentiated from errors via the `status_code` label. Errors
    are identified as any time series with the label `status_code = "STATUS_CODE_ERROR"`.
- `[namespace_]duration_[units]`
  - **Type**: histogram
  - **Description**: a histogram of span durations/latencies. Under the hood, Prometheus histograms
    will create a number of time series. For illustrative purposes, assume no namespace
    is configured and the units are `milliseconds`:
    - `duration_milliseconds_count`: The total number of data points across all buckets in the histogram.
    - `duration_milliseconds_sum`: The sum of all data point values.
    - `duration_milliseconds_bucket`: A collection of `n` time series (where `n` is the number of
      duration buckets) for each duration bucket identified by an `le` (less than
      or equal to) label. The `duration_milliseconds_bucket` counter with lowest `le` and
      `le >= span duration` will be incremented for each span.

The following formula aims to provide some guidance on the number of new time series created:
```
num_status_codes * num_span_kinds * (1 + num_latency_buckets) * num_operations

Where:
  num_status_codes = 3 max (typically 2: ok/error)
  num_span_kinds = 6 max (typically 2: client/server)
  num_latency_buckets = 17 default
```

Plugging those numbers in, assuming default configuration:
```
max = 324 * num_operations
typical = 72 * num_operations
```

Note:
- Custom [duration buckets][spanmetrics-config-duration] or [dimensions][spanmetrics-config-dimensions]
  configured in the spanmetrics connector will alter the calculation above.
- Querying custom dimensions are not supported by SPM and will be aggregated over.

## API

### gRPC/Protobuf

The recommended way to programmatically retrieve RED metrics is via `jaeger.api_v2.metrics.MetricsQueryService` gRPC endpoint defined in the [metricsquery.proto][metricsquery.proto] IDL file.

### HTTP JSON

Used internally by the Monitor tab of Jaeger UI to populate the metrics for its visualizations.

Refer to [this README file][http-api-readme] for a detailed specification of
the HTTP API.

## Troubleshooting

### Check Jaeger-Prometheus connectivity

Verify that Jaeger *query** can connect to Prometheus-compatible metric store by inspecting Jaeger's internal telemetry.

The Jaeger configuration needs to have a metrics endpoint enabled in the `telemetry:` section. Note that the internal telemetry should be exposed on a different port (e.g. `8888`) than the port used to export metrics from the `spanmetrics` connector (e.g. `8889`).

```yaml
service:
  ...
  telemetry:
    resource:
      service.name: jaeger
    metrics:
      level: detailed
      address: 0.0.0.0:8888
```

The `/metrics` endpoint on this port can be used to check if UI queries for SPM data are successful:

```shell
curl -s http://jaeger:8888/metrics | grep jaeger_metricstore
```

The following metrics are of most interest:
  * `jaeger_metricstore_requests_total`
  * `jaeger_metricstore_latency_bucket`

Each of these metrics will have a label for each of the following operations:
  * `get_call_rates`
  * `get_error_rates`
  * `get_latencies`
  * `get_min_step_duration`

If things are working as expected, the metrics with label `result="ok"` should
be incrementing, and `result="err"` being static. For example:
```shell
jaeger_metricstore_requests_total{operation="get_call_rates",result="ok"} 18
jaeger_metricstore_requests_total{operation="get_error_rates",result="ok"} 18
jaeger_metricstore_requests_total{operation="get_latencies",result="ok"} 36

jaeger_metricstore_latency_bucket{operation="get_call_rates",result="ok",le="0.005"} 5
jaeger_metricstore_latency_bucket{operation="get_call_rates",result="ok",le="0.01"} 13
jaeger_metricstore_latency_bucket{operation="get_call_rates",result="ok",le="0.025"} 18

jaeger_metricstore_latency_bucket{operation="get_error_rates",result="ok",le="0.005"} 7
jaeger_metricstore_latency_bucket{operation="get_error_rates",result="ok",le="0.01"} 13
jaeger_metricstore_latency_bucket{operation="get_error_rates",result="ok",le="0.025"} 18

jaeger_metricstore_latency_bucket{operation="get_latencies",result="ok",le="0.005"} 7
jaeger_metricstore_latency_bucket{operation="get_latencies",result="ok",le="0.01"} 25
jaeger_metricstore_latency_bucket{operation="get_latencies",result="ok",le="0.025"} 36
```

If there are issues reading metrics from Prometheus such as a failure to reach
the Prometheus server, then the `result="err"` metrics will be incremented. For example:
```shell
jaeger_metricstore_requests_total{operation="get_call_rates",result="err"} 4
jaeger_metricstore_requests_total{operation="get_error_rates",result="err"} 4
jaeger_metricstore_requests_total{operation="get_latencies",result="err"} 8
```

At this point, checking the logs will provide more insight towards root causing
the problem.

### Query Prometheus

Graphs may still appear empty even when the above Jaeger metrics indicate successful reads
from Prometheus. In this case, query Prometheus directly on any of the metrics that should be generated by the `spanmetrics` connector:

- `traces_span_metrics_duration_milliseconds_bucket`
- `traces_span_metrics_calls_total`

You should expect to see these counters increasing as traces are being received by Jaeger.

### Check the Logs

If the above metrics are present in Prometheus, but not appearing in the Monitor
tab, it means there is a discrepancy between what metrics Jaeger expects to see in
Prometheus and what metrics are actually available.

This can be confirmed by increasing the log level:

```yaml
service:
  telemetry:
    ...
    logs:
      level: debug
```

Outputting logs that resemble the following (formatted for readability):
```
2024-11-26T19:09:43.152Z debug metricsstore/reader.go:258 Prometheus query results
{
  "kind": "extension",
  "name": "jaeger_storage",
  "results": "",
  "query": "sum(rate(traces_span_metrics_calls_total{service_name =~ \"redis\", span_kind =~ \"SPAN_KIND_SERVER\"}[10m])) by (service_name,span_name)",
  "range": {
    "Start": "2024-11-26T19:04:43.14Z",
    "End": "2024-11-26T19:09:43.14Z",
    "Step": 60000000000
  }
}
```

In this instance, let's say OpenTelemetry Collector's `prometheusexporter` introduced
a breaking change that appends a `_total` suffix to counter metrics and the duration units within
histogram metrics (e.g. `duration_milliseconds_bucket`). As we discovered,
Jaeger is looking for the `calls` (and `duration_bucket`) metric names,
while the OpenTelemetry Collector is writing `calls_total` (and `duration_milliseconds_bucket`).

The resolution, in this specific case, is to pass parameters to the metrics backend configuration telling Jaeger
to normalize the metric names such that it knows to search for `calls_total` and
`duration_milliseconds_bucket` instead, like so:

```shell
extensions:
  jaeger_storage:
    backends:
      ...
    metric_backends:
      some_metrics_storage:
        prometheus:
          endpoint: http://prometheus:9090
          normalize_calls: true
          normalize_duration: true
```

### Checking OpenTelemetry Collector Config

If there are error spans appearing in Jaeger, but no corresponding error metrics:

- Check that raw metrics in Prometheus generated by the spanmetrics connector
  (as listed above: `calls`, `calls_total`, `duration_bucket`, etc.) contain
  the `status.code` label in the metric that the span should belong to.
- If there are no `status.code` labels, check the OpenTelemetry Collector
  configuration file, particularly for the presence of the following configuration:
  ```yaml
  exclude_dimensions: ['status.code']
  ```
  This label is used by Jaeger to determine if a request is erroneous.

### Inspect the OpenTelemetry Collector

If the above `latency_bucket` and `calls_total` metrics are empty, then it could
be misconfiguration in the OpenTelemetry Collector or anything upstream from it.

Some questions to ask while troubleshooting are:
- Is the OpenTelemetry Collector configured correctly?
  - See: https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/connector/spanmetricsconnector
- Is the Prometheus server reachable by the OpenTelemetry Collector?
- Are the services sending spans to the OpenTelemetry Collector?
  - See: https://opentelemetry.io/docs/collector/troubleshooting/

### Service/Operation missing in Monitor Tab

If the service/operation is missing in the Monitor Tab, but visible in the Jaeger
Trace search service and operation drop-downs menus, a common cause of this is
the default `server` span kind used in metrics queries.

The service/operations you are not seeing could be from spans that are non-server
span kinds such as client or worse, `unspecified`. Hence, this is an instrumentation
data quality issue, and the instrumentation should set the span kind.

The reason for defaulting to `server` span kinds is to avoid double-counting
both ingress and egress spans in the `server` and `client` span kinds, respectively.

[spm-demo]: https://github.com/jaegertracing/jaeger/tree/main/docker-compose/monitor
[metricsquery.proto]: https://github.com/jaegertracing/jaeger/blob/main/model/proto/metrics/metricsquery.proto
[openmetrics.proto]: https://github.com/jaegertracing/jaeger/blob/main/model/proto/metrics/openmetrics.proto#L53
[opentelemetry-collector]: https://opentelemetry.io/docs/collector/
[spanmetrics-conn]: https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/connector/spanmetricsconnector/README.md
[prom-metric-labels]: https://prometheus.io/docs/concepts/data_model/#metric-names-and-labels
[http-api-readme]: https://github.com/jaegertracing/jaeger/tree/main/docker-compose/monitor#http-api
[spanmetrics-config-dimensions]: https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/connector/spanmetricsconnector/testdata/config.yaml#L23
[spanmetrics-config-duration]: https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/connector/spanmetricsconnector/testdata/config.yaml#L14

### 403 when executing metrics query

If logs contain the error resembling: `failed executing metrics query: client_error: client error: 403`,
it is possible that the Prometheus server is expecting a bearer token.

Jaeger can be configured to pass the bearer token in the metrics queries. The token can be defined via the `token_file_path:` property:
```yaml
extensions:
  jaeger_storage:
    backends:
      ...
    metric_backends:
      some_metrics_storage:
        prometheus:
          endpoint: http://prometheus:9090
          token_file_path: /path/to/token/file
          token_override_from_context: true
```
