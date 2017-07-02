StatsD + Graphite + Grafana Dashboard (SGG Stack)
-------------------------------------------------

* [Grafana](http://metrics.local/)
* [Graphite](http://metrics.local:8000/)
* [ElasticSearch](http://metrics.local/elasticsearch/)

## Creating Metrics

Statsd makes this really simple, to test connectivity and that metrics are flowing through, issue the following command to send a counter with the name 'mattisawesome' to statsd:

* `echo "mattisawesome:1000|c" | nc -w0 -u statsd.local 8125` (GNU/Linux)
* `echo "mattisawesome:1000|c" | nc   -c -u statsd.local 8125` (BSD/Mac OSX)

You should then be able to browse to `http://statsd.local:8000/` and expand `stats > counters > mattisawesome` to see the counter updating.

StatsD + Graphite + Grafana + Kamon Dashboard
---------------------------------------------

This image contains a sensible default configuration of StatsD, Graphite and Grafana, and comes bundled with a example
dashboard that gives you the basic metrics currently collected by Kamon for both Actors and Traces. There are two ways
for using this image:

The container exposes the following ports:

- `80`: Grafana web interface.
- `2003`: Carbon port.
- `8000`: Graphite web interface.
- `8125`: StatsD port.
- `8126`: StatsD administrative port.
- `9200`: ElasticSearch API port.
- `9300`: ElasticSearch administrative port.

To start a container with this image you just need to run the following command:

```bash
docker run -d -p 80:80 -p 8125:8125/udp -p 8126:8126 --name grafana-dashboard sgg
```

If you already have services running on your host that are using any of these ports, you may wish to map the container
ports to whatever you want by changing left side number in the `-p` parameters. Find more details about mapping ports
in the [Docker documentation](http://docs.docker.io/use/port_redirection/#port-redirection).


### Building the image yourself ###

```bash
./build.sh
```

### Using the Dashboard ###

Once your container is running all you need to do is open your browser pointing to the host/port you just published and
play with the dashboard at your wish. We hope that you have a lot of fun with this image and that it serves it's
purpose of making your life easier. This should give you an idea of how the dashboard looks like when receiving data
from one of our toy applications:

![Kamon Dashboard](http://kamon.io/assets/img/kamon-statsd-grafana.png)
![System Metrics Dashboard](http://kamon.io/assets/img/kamon-system-metrics.png)
