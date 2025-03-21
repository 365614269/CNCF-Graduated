# RPC

HTTP RPC enables you to manage your Fluentd instance through HTTP endpoints. You can use this feature as a replacement of [Unix signals](signals.md).

It is especially useful for environments where signals are not supported well e.g. Windows. This requires Fluentd to start not with `--no-supervisor` command-line option.

## Configuration

By default, HTTP RPC is not enabled. To use this feature, set the `rpc_endpoint` option:

```text
<system>
  rpc_endpoint 127.0.0.1:24444
</system>
```

Now, you can manage your Fluentd instance using an HTTP client:

```text
$ curl http://127.0.0.1:24444/api/plugins.flushBuffers
{"ok":true}
```

As evident from the output above, each endpoint returns a JSON object as its response.

## HTTP Endpoints

| Endpoint                                    | Replacement of                                                            | Description                                                     | Version |
| :---                                        | :---:                                                                     | :---:                                                           | :---:   |
| `/api/processes.interruptWorkers`           | [SIGINT](signals.md#sigint-or-sigterm)                                    | Stops the daemon.                                               | v1.0    |
| `/api/processes.killWorkers`                | [SIGTERM](signals.md#sigint-or-sigterm)                                   | Stops the daemon.                                               | v1.0    |
| `/api/processes.zeroDowntimeRestart`        | [SIGUSR2](signals.md#sigusr2)                                             | Restarts Fluentd with zero-downtime. (Not supported on Windows) | v1.18   |
| `/api/processes.flushBuffersAndKillWorkers` | [SIGUSR1](signals.md#sigusr1) and [SIGTERM](signals.md#sigint-or-sigterm) | Flushes buffer and stops the daemon.                            | v1.0    |
| `/api/plugins.flushBuffers`                 | [SIGUSR1](signals.md#sigusr1)                                             | Flushes the buffered messages.                                  | v1.0    |
| `/api/config.reload`                        | [SIGHUP](signals.md#sighup)                                               | Reloads configuration.                                          | v1.0    |
| `/api/config.gracefulReload`                | ---                                                                       | Reloads configuration.                                          | v1.9    |

Appendix:

* `/api/config.gracefulReload`: This is the replacement of `SIGUSR2` before v1.18. Please use `/api/processes.zeroDowntimeRestart` or `/api/config.reload` unless there is a special reason. See [SIGUSR2](signals.md#sigusr2) for details.

If this article is incorrect or outdated, or omits critical information, please [let us know](https://github.com/fluent/fluentd-docs-gitbook/issues?state=open). [Fluentd](http://www.fluentd.org/) is an open-source project under [Cloud Native Computing Foundation \(CNCF\)](https://cncf.io/). All components are available under the Apache 2 License.

