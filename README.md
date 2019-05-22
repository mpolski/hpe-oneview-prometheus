# hpe-oneview-exporter

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

### Option A - build and run binary

 1. Clone repo:

```
$ git clone https://github.com/mpolski/hpe-oneview-msmb-prometheus
$ cd hpe-oneview-exporter
```
 2. Get all required dependencies: 

```
go get -v -d -tags 'static netgo' github.com/mpolski/oneview-golang/ov
go get -v -d -tags 'static netgo' github.com/prometheus/client_golang/prometheus/promhttp
go get -v -d -tags 'static netgo' github.com/prometheus/client_golang/prometheus

```
 
 3. Build the binary:

```
go build hpe-oneview-exporter.go
```

 4. Set all required env variables and start the binary:

```
OV_ENDPOINT="https://OneView_IP" OV_USERNAME="user" OV_PASSWORD="password" OV_AUTHLOGINDOMAIN="my_domain" ./hpe-oneview-msmb-prometheus
```

### Option B - deploy in a Docker container having Prometheus and Grafana already running in your environment

 1. Clone repo:

```
$ git clone https://github.com/mpolski/hpe-oneview-msmb-prometheus
$ cd hpe-oneview-exporter
```

 2. Deploy in [Docker](https://docker.com/) container using included `Dockerfile` to build image:

```
docker build -t hpe-oneview-exporter:<version> -t hpe-oneview-exporter:latest .
```

 3. Start container after providing the env varialbes in the .env_sample file):
    
```
docker run --detach --restart always --rm \
--publish 8080:8080 \
--env-file .env_sample \
--name hpe-oneview-exporter hpe-oneview-exporter:latest
```
 4. Update your Prometheus configuration by adding the new target:

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

 5. Import [Grafana.com](https://grafana.com) Dashboard ID [10233](https://grafana.com/dashboards/10233) into your Grafana instance.


### Option C - Deploy all 3 projects (HPE OneView exporter, Prometheus and Grafana on a Docker host (with docker-compose)
In case you do not have Prometheus nor Grafana working in your environment, follow these steps to build the whole solution in a few minutes.

 1. Clone the repository as descibed above then build the container
  
```
docker build -t hpe-oneview-exporter:<version> -t hpe-oneview-exporter:latest .
```
 2. Add your env variables as per .env_sample file

 3. Start all three containers with docker-compose
  
```
docker-compose -p exporter up -d
```
 4. Verify all three containers are running

```
docker-compose -p exporter ps
```
 5. Browse to Prometheus UI ```http://<node-ip-addr>:9090```, then Status/Targets. Verify the target shows as up.

 6. Browse to Grafana UI ```http://<node-ip-addr>:3000```, first time logon credentials are username: admin, password: admin. 
    The data source should be loaded, import [Grafana.com](https://grafana.com) Dashboard ID [10233](https://grafana.com/dashboards/10233) into your Grafana instance.
  