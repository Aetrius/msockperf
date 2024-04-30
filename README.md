
# Monitoring + SockPerf 
Use at your own discretion.

```                                                                    
```
## Build Process for binary

  `bash
    make build
  ` 

## TODO
- Add env vars into the golang code, or make it as a param to pass in. This will allow targeting the server host and startup params.

## Running Test Locally

  `bash
    docker-compose build -f docker-compose.yml
    docker-compose up -f docker-compose-server.yml # Target whatever host you want to run from as a server.
    docker-compose up -f docker-compose.yml # Applies to the target host you want as a client. Keep in mind I haven't added the env var yet.
  `

## Run
Run docker-compose to run this exporter adhoc OR run the go code.


## HELM INSTALLATION / UPGRADE
helm upgrade msockperf infra/k8s/base/app/msockperf-chart -f infra/k8s/base/app/msockperf-chart/env/values.local.yaml --namespace sre-msockperf-exporter



### Example Image
![alt text](Images/sockperf.png)

### Metric Examples

` ` 