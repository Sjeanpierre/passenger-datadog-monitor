# passenger-datadog-monitor
Send health metrics from Phusion Passenger to DataDog using the StatD interface on the server agent.

#### Purpose
Graph and track Passenger threads and possibly detect misbehaving threads before they become a problem.

###Tracked Metrics
* Processed requests: min,max,average,total
* Memory usage: min,max,average,total
* Thread uptime: min, max, average
* Request queue depth
* Threads in use *vs* Max thread configured


#### Usage
Runs as a daemon with a 10 second sampling interval. Monit, God, SupervisorD, or any other daemon management tool should be used to manage the process.

Sample Monit config

```
check process passenger-datadog-monitor with pidfile /var/run/passenger-datadog-monitor.pid
start program = "/etc/init.d/passenger-datadog-monitor start"
stop  program = "/etc/init.d/passenger-datadog-monitor stop"
```

`./passenger-datadog-monitor` as root, since access to passenger-status requires root.

`./passenger-datadig-monitor print` for basic console output of stats, useful for debugging

[udp.rb](https://github.com/Sjeanpierre/passenger-datadog-monitor/blob/master/udp.rb) can be run locally when you want to see what is being recieved on the server side.

