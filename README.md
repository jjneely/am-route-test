Alertmanager Route Test
=======================

I use Alertmanager[1] with Prometheus to route alerts from Prometheus servers
to the right team.  This also does various filtering, handle high and low
priority alerts, and other things which have made the routing tree more complex
than I'd like.  Often, requests to fine tune various bits come in that require
some modification of the route tree.  Because alert notification simply cannot
break, I wanted a tool that could automate unit testing of my route tree.

This lead to writing of `am-route-test`.  This will parse a given Alertmanager
config, use the label / value pairs supplied on the command line, and print
out the Alertmanager receivers that would eventually handle that an alert
with those labels.

```
Usage of am-route-test:
Evalute what receiver(s) are used when an alert with the given labels
is passed through the routing tree found in the Alertmanager config.

        am-route-test [options] label=value [label=value ...]

  -e string
        Comma seperated list of expected alert receivers
  -f string
        Alertmanager configuration file (default "alertmanager.yml")

```

Use the `-e` option to make the tool exit with a status code of 2 if the
receivers found do not match the comma seperated list of receivers passed to
that flag.

Examples are below for some edge case testing we have in the PERF environment.

```
#!/bin/bash

set +eux

# ...

# PERF environment is special
./am-route-test -e perfcap-high-priority,default-event-handler \
    monitor=node severity=page environment=perf
./am-route-test -e perfcap-low-priority,default-event-handler \
    monitor=node severity=warning environment=perf
```

Maintenance
-----------

Alertmanager is not meant to be used as a library, so its completely vendored
into this source code.  Maintenance and upgrades should include updating
the vendoring.  Perhaps even managing it with `govendor` or `dep`.

Could this be contributed directly to the `amtool` CLI that ships with
Alertmanager?

Contributing
------------

Open a PR.

[1]: https://github.com/prometheus/alertmanager
