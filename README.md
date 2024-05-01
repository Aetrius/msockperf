
# Monitoring + SockPerf 
The code is provided as-is with no warranties.

## Usage

[Helm](https://helm.sh) must be installed to use the charts.
Please refer to Helm's [documentation](https://helm.sh/docs/) to get started.

Once Helm is set up properly, add the repo as follows:

```console
helm repo add aetrius/msockperf https://aetrius.github.io/msockperf

helm install msockperf

or 

helm install msockperf infra/k8s/base/app/msockperf-chart -f infra/k8s/base/app/msockperf-chart/env/values.local.yaml

helm upgrade msockperf infra/k8s/base/app/msockperf-chart -f infra/k8s/base/app/msockperf-chart/env/values.local.yaml
```



You can then run `helm search repo aetrius/msockperf` to see the charts.

<!-- Keep full URL links to repo files because this README syncs from main to gh-pages.  -->
Chart documentation is available in [aetrius/msockperf directory](https://github.com/aetrius/msockperf/README.md).


## Values
Values table

| Parameter            | Description                                 | Default        | Possible Values |
|----------------------|---------------------------------------------|----------------|-----------------|
| `namespace.enabled` | Enable creation of namespace                | `true`         | `true` or `false` |
| `namespace.name`     | Namespace name                              | `sre-msockperf-exporter` | String (namespace name) |
| `environmentDomain`  | Domain for environment                       | `msockperf.domainname.com` | String (domain name) |
| `labels.enabled`     | Enable labels                                | `true`         | `true` or `false` |
| `labels.owner`       | Owner label                                  | `sre`          | String (owner name) |
| `ingress.enabled`    | Enable ingress                               | `true`         | `true` or `false` |
| `servicemonitor.enabled` | Enable Prometheus ServiceMonitor          | `true`         | `true` or `false` |
| `portworx.enabled`   | Enable Portworx storage                     | `false`        | `true` or `false` |
| `rancher.enabled`    | Enable Rancher integration                  | `true`         | `true` or `false` |
| `rancher.projectID`  | Rancher project ID                          | (empty)        | String (project ID) |
| `msockperf.serviceName` | Name of Msockperf service                | `msockperf-service` | String (service name) |
| `msockperf.appName`  | Name of Msockperf app                       | `msockperf-app` | String (app name) |
| `msockperf.image`    | Docker image for Msockperf                  | `ghcr.io/aetrius/msockperf-client/msockperf-client:main` | String (image URL) |
| `msockperf.replicaCount` | Number of Msockperf replicas            | `1`            | Integer (replica count) |
| `msockperf.revisionHistoryLimit` | Revision history limit                | `0`            | Integer (revision limit) |
| `msockperf.webPort`  | Port for Msockperf web interface           | `8082`         | Integer (port number) |
| `msockperf.namespace` | Namespace for Msockperf                   | `default`      | String (namespace name) |
| `msockperf.remoteHost` | Remote host for Msockperf                | `sockperf-server-service.sre-msockperf-exporter` | String (host) |
| `msockperf.remotePort` | Remote port for Msockperf                | `3600`         | Integer (port number) |
| `msockperf.servicePort` | Service port for Msockperf              | `8000`         | Integer (port number) |
| `msockperf.imagePullPolicy` | Image pull policy for Msockperf       | `Always`       | `Always`, `IfNotPresent`, or `Never` |
| `sockperfServer.serviceName` | Name of Sockperf Server service        | `sockperf-server-service` | String (service name) |
| `sockperfServer.appName` | Name of Sockperf Server app             | `sockperf-app` | String (app name) |
| `sockperfServer.image` | Docker image for Sockperf Server         | `ghcr.io/aetrius/msockperf-server/msockperf-server:main` | String (image URL) |
| `sockperfServer.replicaCount` | Number of Sockperf Server replicas    | `1`            | Integer (replica count) |
| `sockperfServer.revisionHistoryLimit` | Revision history limit           | `1`            | Integer (revision limit) |
| `sockperfServer.sockperfPort` | Port for Sockperf Server               | `3600`         | Integer (port number) |
| `sockperfServer.webPort` | Web port for Sockperf Server             | `3600`         | Integer (port number) |


## Build Process for binary locally

  `bash
    make build
  ` 

## Contributing

<!-- Keep full URL links to repo files because this README syncs from main to gh-pages.  -->
We'd love to have you contribute! Please refer to our [contribution guidelines](https://github.com/aetrius/msockperf/CONTRIBUTING.md) for details.


## Running Test Locally in Docker

  `bash

    docker-compose build -f docker-compose.yml

    docker-compose up -f docker-compose-server.yml # Target whatever host you want to run from as a server.

    docker-compose up -f docker-compose.yml # Applies to the target host you want as a client. Keep in mind I haven't added the env var yet.
    
  `


## License

<!-- Keep full URL links to repo files because this README syncs from main to gh-pages.  -->
[Apache 2.0 License](https://github.com/aetrius/msockperf/LICENSE).


