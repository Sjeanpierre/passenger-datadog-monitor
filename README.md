# passenger-datadog-monitor

Send health metrics from Phusion Passenger to DataDog using the StatsD interface on the server agent.

## Purpose

Graph and track Passenger threads and possibly detect misbehaving threads before they become a problem.

## Tracked Metrics

### Aggregated

- Processed requests: min,max,average,total
- Memory usage: min,max,average,total
- Thread uptime: min, max, average
- Request queue depth
- Threads in use _vs_ Max thread configured & Threads used between runs

### Discrete

- Process memory usage per Passenger Process PID
- OS thread count per Passenger Process PID
- Requests processed per Passenger Process PID
- Process Idle time per Passenger Process PID

## Installation

### Downloading from Github

The `passenger-datadog-monitor` binary can be downloaded from the [releases area](https://github.com/Sjeanpierre/passenger-datadog-monitor/releases) of this repository for Linux

### Building the binary

You will first need to build the `passenger-datadog-monitor` executable using [Go](https://golang.org). You can download the source and dependencies, and build the binary by running:

```sh
go get -v github.com/Sjeanpierre/passenger-datadog-monitor
```

Once it completes, you should find your new `passenger-datadog-monitor` executable in your `$GOROOT/bin` directory.

Note that if you are building in a different environment from where you plan to deploy, you should configure your [target operating system and architecture](https://golang.org/doc/install/source#environment).

The [Makefile](Makefile) in this repository will cross compile for Linux.

### Installing the binary

After you've built the executable, you should install it on your server (e.g. in `/usr/bin/`).

## Usage

`passenger-datadog-monitor` runs as a daemon with a 10 second sampling interval. Monit, God, SupervisorD, or any other daemon management tool should be used to manage the process.

Sample Monit config

```plaintext
check process passenger-datadog-monitor with pidfile /var/run/passenger-datadog-monitor.pid
start program = "/etc/init.d/passenger-datadog-monitor start"
stop  program = "/etc/init.d/passenger-datadog-monitor stop"
```

You should run `passenger-datadog-monitor` as root, since access to passenger-status requires root.

### Flags

| flag | type | description | example |
|:-----|:---|:------------|:---|
| -host | string | StatsD collector IP - useful when running with a Kubernetes DaemonSet or other external collector | -host=100.124.102.21 |
| -port | string | StatsD collector UDP Port - useful when running in Docker or other custom environments | -port=81333 |
| -print | bool | Enable debug and stats printing | -print |

Full example:

```sh
passenger-datadog-monitor -host=$STATSD_IP -port=$STATSD_PORT
```

### Testing

[udp.rb](https://github.com/Sjeanpierre/passenger-datadog-monitor/blob/master/server/udp.rb) can be run locally when you want to see what is being received on the server side.

Alternatively you can listen using netcat: `nc -kulvw 0 8125`
