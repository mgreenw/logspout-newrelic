# Newrelic LogSpout
A Docker image to stream logs from your containers to Newrelic. Based on [LogDNA Logspout](https://github.com/logdna/logspout).

NOTE: This is not production ready.

## How to Use

### Environment Variables
The following variables can be used to tune `LogSpout` for specific use cases.

#### Log Router Specific
* __FILTER_NAME__: Filter by Container Name with Wildcards, *Optional*
* __FILTER_ID__: Filter by Container ID with Wildcards, *Optional*
* __FILTER_SOURCES__: Filter by Comma-Separated List of Sources, *Optional*
* __FILTER_LABELS__: Filter by Comma-Separated List of Labels, *Optional*

__Note__: More information can be found [here](https://github.com/gliderlabs/logspout/tree/0da75a223db992cd5abc836796174588ddfc62b4/routesapi#routes-resource).

#### Ingestion Specific
* __NEWRELIC_KEY__: Newrelic License Key, *Required*
* __HOSTNAME__: Alternative Hostname, *Optional*
  * __Default__: System's Hostname
* __NEWRELIC_URL__: Specific Endpoint to Stream Log into, *Optional*
  * __Default__: `log-api.newrelic.com/log/v1`

__Note__: Logging the `LogSpout` Container is recommended to keep track of HTTP Request Errors or Exceptions.

#### Limits
* __FLUSH_INTERVAL__: How frequently batches of logs are sent (in `milliseconds`), *Optional*
  * __Default__: 250
* __HTTP_CLIENT_TIMEOUT__: Time limit (in `seconds`) for requests made by this HTTP Client, *Optional*
  * __Default__: 30
  * __Source__: [net/http/client.go#Timeout](https://github.com/golang/go/blob/master/src/net/http/client.go#L89-L104)
* __INACTIVITY_TIMEOUT__: How long to wait for inactivity before declaring failure in the `Docker API` and restarting, *Optional*
  * __Default__: 1m
  * __Note__: More information about the possible values can be found [here](https://github.com/gliderlabs/logspout#detecting-timeouts-in-docker-log-streams). Also see [`time.ParseDuration`](https://golang.org/pkg/time/#ParseDuration) for valid format as recommended [here](https://github.com/gliderlabs/logspout/blob/e671009d9df10e8139f6a4bea8adc9c7878ff4e9/router/pump.go#L112-L116).
* __MAX_BUFFER_SIZE__: The maximum size (in `mb`) of batches to ship to `Newrelic`, *Optional*
  * __Default__: 2
* __MAX_REQUEST_RETRY__: The maximum number of retries for sending a line when there are network failures, *Optional*
  * __Default__: 5

### Docker
Create and run container named *logspout* from this image using CLI:
```bash
sudo docker run --name="logspout-newrelic" --restart=always \
-d -v=/var/run/docker.sock:/var/run/docker.sock \
-e NEWRELIC_KEY="<Newrelic License Key>" \
mgreenw/logspout-newrelic:latest
```

### Docker Cloud
Append the following to your Docker Cloud stackfile:
```yaml
logspout-newrelic:
  autoredeploy: true
  deployment_strategy: every_node
  environment:
    - NEWRELIC_KEY="<Newrelic License Key Key>"
  image: mgreenw/logspout-newrelic:latest
  restart: always
  volumes:
    - '/var/run/docker.sock:/var/run/docker.sock'
```
