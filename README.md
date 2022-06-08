# openconnect vyos telemetry (ovt)

OVT sends openconnect telemetry to InfluxDB (v1) and Loki (TODO).

While this project is targeted at VyOS, it can run on any server executing an
openconnect instance.

OVT will send metrics (bytes/connection age) to InfluxDB and events
(connect/disconnect) to Loki.

## Usage

```
Usage:
  ovt [OPTIONS]

Application Options:
  -L, --loki.url=      Loki server URL
  -U, --loki.user=     Loki basic auth username
  -P, --loki.pass=     Loki basic auth password
  -O, --loki.orgid=    Loki X-Scope-OrgID header to add
  -I, --influx.url=    InfluxDB server URL
  -V, --influx.user=   InfluxDB basic auth username
  -Q, --influx.pass=   InfluxDB basic auth password
  -D, --influx.db=     InfluxDB database name
  -d, --metrics.delay= Delay between metrics collection in seconds (default: 2)
  -M, --metrics.cmd=   Command to run to fetch metrics (default: sudo occtl -s /run/ocserv/occtl.socket -j show users)
  -E, --events.cmd=    Command to run to fetch events (default: sudo occtl -s /run/ocserv/occtl.socket show events)
      --debug          Enable debug logging
  -n, --dry-run        Dry run, do not send any data to InfluxDB
  -v, --version        displays versions

Help Options:
  -h, --help           Show this help message
```


