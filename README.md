# passenger-datadog-monitor
Send health metrics from Phusion Passenger to DataDog using the StatsD interface on the server agent.

## Purpose
Graph and track Passenger threads and possibly detect misbehaving threads before they become a problem.

## Tracked Metrics
* Processed requests: min,max,average,total
* Memory usage: min,max,average,total
* Thread uptime: min, max, average
* Request queue depth
* Threads in use *vs* Max thread configured

## Installation
### Building the binary
You will first need to build the `passenger-datadog-monitor` executable using [Go](https://golang.org). Download the source and dependencies, then build the binary running:
```
go get -v github.com/Sjeanpierre/passenger-datadog-monitor
```
Once it completes, you should find your new `passenger-datadog-monitor` executable in your `$GOROOT/bin` directory.

Note that if you are building in a different environment from where you plan to deploy, you should configure your [target operating system and architecture](https://golang.org/doc/install/source#environment).

### Installing the binary
After you've built the executable, you should install it on your server (e.g. in `/usr/bin/`).

## Usage
`passenger-datadog-monitor` runs as a daemon with a 10 second sampling interval. Monit, God, SupervisorD, or any other daemon management tool should be used to manage the process.

Sample Monit config

```
check process passenger-datadog-monitor with pidfile /var/run/passenger-datadog-monitor.pid
start program = "/etc/init.d/passenger-datadog-monitor start"
stop  program = "/etc/init.d/passenger-datadog-monitor stop"
```

`./passenger-datadog-monitor` as root, since access to passenger-status requires root.

`./passenger-datadog-monitor print` for basic console output of stats, useful for debugging

[udp.rb](https://github.com/Sjeanpierre/passenger-datadog-monitor/blob/master/server/udp.rb) can be run locally when you want to see what is being received on the server side.
