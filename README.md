# passenger-datadog-monitor
Send health metrics from Phusion Passenger to DataDog agent using the StatD interface on the.

#### Purpose
Graph and track Passenger threads and possibly detect misbehaving threads before they become a problem.

###Tracked Metrics
* Processed requests: min,max,average,total
* Memory usage: min,max,average,total
* Thread uptime: min, max, average
* Request queue depth
* Threads in use *vs* Max thread configured


#### Usage
I currently use this in a cron task that runs every 10 or so seconds. a loop can also be added to the main function to achieve similar results.

`./passenger-datadog-monitor` as root, since access to passenger-status requires root.

`./passenger-datadig-monitor print` for basic console output of stats, useful for debugging

[udp.rb](https://github.com/Sjeanpierre/passenger-datadog-monitor/blob/master/udp.rb) can be run locally when you want to see what is being recieved on the server side.

