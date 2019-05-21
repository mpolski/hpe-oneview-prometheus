# hpe-oneview-prometheus

## HPE OneView exporter for Prometheus

This module provides an HTTP endpoint for [Prometheus](https://prometheus.io/) time series database and alerting system. Prometheus to will scrap metrics for all server-hardware and enclosures managed by a given [HPE OneView](https://www.hpe.com/pl/en/integrated-systems/software.html) instance.
A single HPE OneView instance manages multiple resources like server-hardware and interconnect modules hosted in enclsoures. HPE OneView REST API has been used to collect status of abovementioned elements as well as the following environment data as reported by the hardware: power consumption (average and peak), ambient temperature, CPU utilization as well as state of components. This data will be exposed at an HTTP endpoint in a Prometheus compatible format. The metrics can be easily visualized using [Grafana](https://grafana.com/) dashboards. A sample dashboard to visualize HPE OneView Prometheus metrics is available at https://grafana.com, dashboard ID 9984.

## Configuration variables

Configuration options set via environment variables:
* `OV_ENDPOINT`: \<your HPE OneView instance IP address\> - required, e.g. "https://10.10.0.50"
* `OV_USERNAME`: \<your HPE OneView username\> - required, e.g. "administrator"
* `OV_PASSWORD`: \<your HPE OneView password\> - required, e.g. "password"
* `OV_AUTHLOGINDOMAIN`: \[\<your HPE OneView login domain\>\] - optional, not used if empty

## HPE OneView exporter for Prometheus - deployment options

To deploy this application, execute the following commands:

  1. Clone repo:

```
$ git clone https://github.com/mpolski/hpe-oneview-msmb-prometheus
$ cd hpe-oneview-prometheus
```

  2. Get all required dependencies: 

```
go get -v -d -tags 'static netgo' github.com/mpolski/oneview-golang/ov
go get -v -d -tags 'static netgo' github.com/prometheus/client_golang/prometheus/promhttp
go get -v -d -tags 'static netgo' github.com/prometheus/client_golang/prometheus

```

### Option A - build and run binary

  1. To build binary:

```
go build hpe-oneview-prometheus.go
```

  2. Then set all required env variables and start binary:

```
OV_ENDPOINT="https://OneView_IP" OV_USERNAME="user" OV_PASSWORD="password" OV_AUTHLOGINDOMAIN="my_domain" ./hpe-oneview-msmb-prometheus
```

### Option B - deploy in a Docker container

In case you run Prometheus and Grafana in your environment, deploy this container first.

  1. Deploy in [Docker](https://docker.com/) container using included `Dockerfile` to build image:

```
docker build -t hpe-oneview-prometheus:<version> -t hpe-oneview-prometheus:latest .
```

  2. Then set all required env variables and start container, ex:
    
```
docker run --detach --restart always --rm \
--publish 8080:8080 \
--env OV_ENDPOINT="OneView_IP" \
--env OV_USERNAME="user" \
--env OV_PASSWORD="password" \
--env OV_AUTHLOGINDOMAIN="my_domain" \
--name hpe-oneview-prometheus hpe-oneview-prometheus:latest
```
Now, update your Prometheus configuration by adding the new target:

```
- job_name: hpe_oneview_exporter
  scrape_interval: 5m
  scrape_timeout: 10s
  metrics_path: /metrics
  scheme: http
  static_configs:
  - targets:
    - <container_IP_or_name>:8080
```

Lastly, import [Grafana.com](https://grafana.com) Dashboard ID [10233](https://grafana.com/dashboards/10233) into your Grafana instance.


### Option C - Prometheus and Grafana Deployment in Docker host (with docker-compose)
In case you do not have Prometheus nor Grafana working in your environment, follow these steps to build the whole solution in a few minutes.

  1. Clone the repository as descibed above then build the container
  
```
docker build -t hpe-oneview-prometheus:<version> -t hpe-oneview-prometheus:latest .
```
  2. Start all three containers with docker-compose
  
```
docker-compose -p exporter up -d
```
  3. Verify all three containers are running
  
```
docker-compose -p exporter ps
```
  4. Browse to Prometheus UI ```http://<node-ip-addr>:9090```, then Status/Targets. Verify the target shows as up.
  5. Browse to Grafana UI ```http://<node-ip-addr>:3000```, (first time logon credentials are username: admin, password: admin). The data source should be loaded, import the Dashboard using id 10233.
  