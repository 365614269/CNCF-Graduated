# Plugin Helper: Compat Parameters

`compat_parameters` helper convert parameters from v0.12 to v1.0 style.

Here is an example:

```ruby
require 'fluent/plugin/output'

module Fluent::Plugin
  class ExampleOutput < Output
    Fluent::Plugin.register_output('example', self)

    # 1. Load `compat_parameters` helper
    helpers :compat_parameters

    # Omit `start`, `shutdown` and other plugin APIs

    def configure(conf)
      compat_parameters_convert(conf, :buffer, :inject, default_chunk_key: 'time')
      super

      # ...
    end
  end
end
```

## Methods

### `compat_parameters_convert(conf, *types, **kwargs)`

* `conf`: `Fluent::Configuration` instance
* `types`: Following types are supported:
  * `:buffer`
  * `:inject`
  * `:extract`
  * `:parser`
  * `:formatter`
* `kwargs`:
  * `default_chunk_key`: Sets the default chunk key. For more details,

    see [Buffer Section Configurations](../configuration/buffer-section.md).

**IMPORTANT**: You cannot mix v1.0 and v0.12 styles in one plugin directive. If you mix v1.0 and v0.12 styles, v1.0 style is used and v0.12 style is ignored.

Here is an example:

```text
<match pattern>
  @type foo
  # This flush_interval is ignored because <buffer> exist. No conversion.
  flush_interval 30s
  # This <buffer> parameters are used
  <buffer>
    @type file
    path /path/to/buffer
    retry_max_times 10
    queue_limit_length 256
  </buffer>
</match>
```

#### `buffer`

| old \(v0.12\) | new \(v1\) | note |
| :--- | :--- | :--- |
| `buffer_type` | `@type` |  |
| `buffer_path` | `path` |  |
| `num_threads` | `flush_thread_count` |  |
| `flush_interval` | `flush_interval` |  |
| `try_flush_interval` | `flush_thread_interval` |  |
| `queued_chunk_flush_interval` | `flush_thread_burst_interval` |  |
| `disable_retry_limit` | `retry_forever` |  |
| `retry_limit` | `retry_max_times` |  |
| `max_retry_wait` | `retry_max_interval` |  |
| `buffer_chunk_limit` | `chunk_limit_size` |  |
| `buffer_queue_limit` | `queue_limit_length` |  |
| `buffer_queue_full_action` | `overflow_action` |  |
| `flush_at_shutdown` | `flush_at_shutdown` |  |
| `time_slice_format` | `timekey` | also set `chunk_key` to `time` |
| `time_slice_wait` | `timekey_wait` |  |
| `timezone` | `timekey_zone` |  |
| `localtime` | `timekey_use_utc` | exclusive with `utc` |
| `utc` | `timekey_use_utc` | exclusive with `localtime` |

This flat configuration:

```text
buffer_type file
buffer_path /path/to/buffer
retry_limit 10
flush_interval 30s
buffer_queue_limit 256
```

converts to:

```text
<buffer>
  @type file
  path /path/to/buffer
  retry_max_times 10
  flush_interval 30s
  queue_limit_length 256
</buffer>
```

For more details, see [Buffer Section Configuration](../configuration/buffer-section.md).

#### `inject`

| old \(v0.12\) | new \(v1\) | note |
| :--- | :--- | :--- |
| `include_time_key` | `time_key` | if `true`, set `time_key` |
| `time_key` | `time_key` |  |
| `time_format` | `time_format` |  |
| `timezone` | `timezone` |  |
| `include_tag_key` | `tag_key` | if `true`, set `tag_key` |
| `tag_key` | `tag_key` |  |
| `localtime` | `localtime` | exclusive with `utc` |
| `utc` | `localtime` | exclusive with `localtime` |

This flat configuration:

```text
include_time_key
time_key event_time
utc
```

converts to:

```text
<inject>
  time_key event_time
  localtime false
</inject>
```

For more details, see [Inject Plugin Helper API](api-plugin-helper-inject.md).

#### `extract`

| old \(v0.12\) | new \(v1\) | note |
| :--- | :--- | :--- |
| `time_key` | `time_key` |  |
| `time_format` | `time_format` |  |
| `timezone` | `timezone` |  |
| `tag_key` | `tag_key` |  |
| `localtime` | `localtime` | exclusive with `utc` |
| `utc` | `localtime` | exclusive with `localtime` |

This flat configuration:

```text
include_time_key
time_key event_time
utc
```

converts to:

```text
<extract>
  time_key event_time
  localtime false
</extract>
```

For more details, see [Extract Plugin Helper API](api-plugin-helper-extract.md)

#### `parser`

| old \(v0.12\) | new \(v1\) | note | plugin |
| :--- | :--- | :--- | :--- |
| `format` | `@type` |  |  |
| `types` | `types` | converted to JSON format |  |
| `types_delimiter` | `types` |  |  |
| `types_label_delimiter` | `types` |  |  |
| `keys` | `keys` |  | `CSVParser`, `TSVParser` \(old `ValuesParser`\) |
| `time_key` | `time_key` |  |  |
| `time_format` | `time_format` |  |  |
| `localtime` | `localtime` | exclusive with `utc` |  |
| `utc` | `localtime` | exclusive with `localtime` |  |
| `delimiter` | `delimiter` |  |  |
| `keep_time_key` | `keep_time_key` |  |  |
| `null_empty_string` | `null_empty_string` |  |  |
| `null_value_pattern` | `null_value_pattern` |  |  |
| `json_parser` | `json_parser` |  | `JSONParser` |
| `label_delimiter` | `label_delimiter` |  | `LabeledTSVParser` |
| `format_firstline` | `format_firstline` |  | `MultilineParser` |
| `message_key` | `message_key` |  | `NoneParser` |
| `with_priority` | `with_priority` |  | `SyslogParser` |
| `message_format` | `message_format` |  | `SyslogParser` |
| `rfc5424_time_format` | `rfc5424_time_format` |  | `SyslogParser` |

This flat configuration:

```text
format /^(?<log_time>[^ ]*) (?<message>)+*/
time_key log_time
time_format %Y%m%d%H%M%S
keep_time_key
```

converts into:

```text
<parse>
  @type regexp
  pattern /^(?<log_time>[^ ]*) (?<message>)+*/
  time_key log_time
  time_format %Y%m%d%H%M%S
  keep_time_key
</parse>
```

For more details, see [Parser Plugin Overview](../parser/) and [Writing Parser Plugins](../plugin-development/api-plugin-parser.md).

#### `formatter`

| old \(v0.12\) | new \(v1\) | note | plugin |
| :--- | :--- | :--- | :--- |
| `format` | `@type` |  |  |
| `delimiter` | `delimiter` |  |  |
| `force_quotes` | `force_quotes` |  | `CSVFormatter` |
| `keys` | `keys` |  | `TSVFormatter` |
| `fields` | `fields` |  | `CSVFormatter` |
| `json_parser` | `json_parser` |  | `JSONFormatter` |
| `label_delimiter` | `label_delimiter` |  | `LabeledTSVFormatter` |
| `output_time` | `output_time` |  | `OutFileFormatter` |
| `output_tag` | `output_tag` |  | `OutFileFormatter` |
| `localtime` | `localtime` | exclusive with `utc` | `OutFileFormatter` |
| `utc` | `utc` | exclusive with `localtime` | `OutFileFormatter` |
| `timezone` | `timezone` |  | `OutFileFormatter` |
| `message_key` | `message_key` |  | `SingleValueFormatter` |
| `add_newline` | `add_newline` |  | `SingleValueFormatter` |
| `output_type` | `output_type` |  | `StdoutFormatter` |

This flat configuration:

```text
format json
```

converts to:

```text
<format>
  @type json
</format>
```

For more details, see [Formatter Plugin Overview](../formatter/) and [Writing Formatter Plugins](../plugin-development/api-plugin-formatter.md).

## Plugins using `compat_parameters`

* [`filter_parser`](../filter/parser.md)
* [`in_exec`](../input/exec.md)
* [`in_http`](../input/http.md)
* [`in_syslog`](../input/syslog.md)
* [`in_tail`](../input/tail.md)
* [`in_tcp`](../input/tcp.md)
* [`in_udp`](../input/udp.md)
* [`out_exec`](../output/exec.md)
* [`out_exec_filter`](../output/exec_filter.md)
* [`out_file`](../output/file.md)
* [`out_forward`](../output/forward.md)
* [`out_stdout`](../output/stdout.md)

If this article is incorrect or outdated, or omits critical information, please [let us know](https://github.com/fluent/fluentd-docs-gitbook/issues?state=open). [Fluentd](http://www.fluentd.org/) is an open-source project under [Cloud Native Computing Foundation \(CNCF\)](https://cncf.io/). All components are available under the Apache 2 License.

