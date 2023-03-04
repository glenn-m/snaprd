# Snaprd

[![Go](https://github.com/glenn-m/snaprd/actions/workflows/go.yml/badge.svg)](https://github.com/glenn-m/snaprd/actions/workflows/go.yml)

Daemon that runs [snapraid](https://www.snapraid.it/) on a schedule and surfaces [Prometheus](https://prometheus.io/) metrics.

### Configuration

| Flag          | Env Var             | Description                                         |
|---------------|---------------------|-----------------------------------------------------|
| --config      | SNAPRD_CONFIG_FILE  | The path to the snaprd config yaml file.            |
| --metricsPort | SNAPRD_METRICS_PORT | The port the Prometheus metrics will be exposed on. |
| --metricsPath | SNAPRD_METRICS_PATH | The path the Prometheus metrics will be exposed on. |

Configuration is done via a YAML file. To see the various options available, check out the [example config file](./snaprd.example.yaml).

### TODO

* Create example systemd config
