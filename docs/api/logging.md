# Logging
Throughout the server there is logging at nearly every point in order to potentially give the deepest or shallowest level of insight into the server's actions.
Using the logging library [logrus](https://github.com/sirupsen/logrus), we are able to log at different levels and in two different formats.
Additional fields can be added to the logs that can provide extra information to the log.
While logrus is not incredibly efficient when compared to other libraries like [zap](https://github.com/uber-go/zap#performance) from Uber, it is good enough.

## The Levels
The 7 supported levels are `trace`, `debug`, `info`, `warn`, `error`, `fatal`, and `panic`.
Each one is for different levels of verbosity, with trace being the highest verbosity and panic being the lowest.
We use trace, debug, info, error, and fatal throughout the server.

### Trace
The trace level is logged for deep insight into what the server is doing.
It should never be used when running for a long time as it generates a lot of logs in a very short period of time.
Trace logs include information about the current values the server is processing and non-important errors.
<br>
Example: `time="2015-03-26T01:27:38-04:00" level=trace msg="Started observing beach" animal=walrus number=8`

### Debug
The debug level is used for in-depth information about the server, but not exactly what the server is doing.
It is used sparingly and usually appears at the end of each method call.
As with the trace level, it also contains some extra fields with extra information.
However, this only provides generic information about the user and the operation.
<br>
Example: `time="2015-03-26T01:27:38-04:00" level=debug msg="Started observing beach" animal=walrus number=8`

### Info
The info level is for messages such as status messages to signify events that are good to know, but don't affect anything.
Here, the info level is used during server startup and shutdown to signify major events.
Unlike the grainy debug and trace levels, the info level only notifies about events that provide significant changes.
<br>
Example: `time="2015-03-26T01:27:38-04:00" level=info msg="Started observing beach" animal=walrus number=8`

### Error
The error level is for messages that should happen as infrequently as possible, but are not server stopping if they do.
In this case, error logs occur when a 500 status code is sent to the user as it signifies something when wrong where it should not have.
The data included with these logs are the error itself and any information that can be used to debug the error.
This is similar to the amount of data that would be included in trace or debug levels.
<br>
Example: `time="2015-03-26T01:27:38-04:00" level=error msg="Started observing beach" animal=walrus number=8`

### Fatal
The fatal level are errors that are unrecoverable and require the server to exit.
These errors should happen rarely, if ever.
If they do for some reason happen, they will generally be during startup or shutdown.
<br>
Example: `time="2015-03-26T01:27:38-04:00" level=fatal msg="Started observing beach" animal=walrus number=8`

## The Formats
There are two formats that logrus supports: `text` and `json`.
The text format is intended to be human readable and viewable in a development environment.
The JSON format is intended for ingress via a log parser like Logstash or Prometheus.

### Text
The text format is a tab-formatted string that is printed out on a single line with each field's key and value separated by an equals sign.
It allows for the quick and simple viewing of logs during development.
They are also able to color the text based on the level if a TTY is attached.
<br><br>
Example:
```
time="2015-03-26T01:27:38-04:00" level=trace msg="Started observing beach" animal=walrus number=8
time="2015-03-26T01:27:38-04:00" level=debug msg="Temperature changes" temperature=-4
time="2015-03-26T01:27:38-04:00" level=info msg="A group of walrus emerges from the ocean" animal=walrus size=10
time="2015-03-26T01:27:38-04:00" level=warning msg="The group's number increased tremendously!" number=122 omg=true
time="2015-03-26T01:27:38-04:00" level=error msg="Some peguins swim by" number=122 omg=true
time="2015-03-26T01:27:38-04:00" level=fatal msg="The ice breaks!" err=&{0x2082280c0 map[animal:orca size:9009] 2015-03-26 01:27:38.441574009 -0400 EDT panic It's over 9000!} number=100 omg=true
time="2015-03-26T01:27:38-04:00" level=panic msg="It's over 9000!" animal=orca size=9009
```

### JSON
The JSON format is a string that is printed out on a single line in standard JSON.
It allows for the processing of the logs by an ingress system such as Logstash or Prometheus and then parsed with a tool like Kibana or Grafana.
These systems are typically used for high volume, production system logging that must deal with logs from multiple services.
<br><br>
Example:
```
{"level":"trace","msg":"Started observing beach","animal":"walrus","size":8,"time":"2014-03-10 19:57:38.562264131 -0400 EDT"}
{"level":"debug","msg":"Temperature changes","number":122,"omg":true,"time":"2014-03-10 19:57:38.562471297 -0400 EDT"}
{"level":"info","msg":"A group of walrus emerges from the ocean","animal":"walrus","size":10,"time":"2014-03-10 19:57:38.562500591 -0400 EDT"}
{"level":"warning","msg":"The group's number increased tremendously!","animal":"walrus","size":9,"time":"2014-03-10 19:57:38.562527896 -0400 EDT"}
{"level":"error","msg":"Some peguins swim by","animal":"walrus","size":9,"time":"2014-03-10 19:57:38.562527896 -0400 EDT"}
{"level":"fatal","msg":"The ice breaks!","number":100,"omg":true,"time":"2014-03-10 19:57:38.562543128 -0400 EDT"}
{"level":"panic","msg":"It's over 9000!","number":100,"omg":true,"time":"2014-03-10 19:57:38.562543128 -0400 EDT"}
```

## Logging Requests
On each request, there is information logged about it after it has completed.
This is done through the hijacking of a response writer to save the length and status code of the response.
While there are request loggers like the logging handler from [gorilla](https://github.com/gorilla/handlers), this one allows us to log in our own format.
The code for hijacking the response writer is actually taken from that same logger implementation.
The logger gives the request URI, remote address of the requester, protocol used, request method, response status, and response size.
As it is using logrus, we are able to log in either human-readable text or machine-readable JSON.
