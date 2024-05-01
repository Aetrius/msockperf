helm install msockperf infra/k8s/base/app/msockperf-chart -f infra/k8s/base/app/msockperf-chart/env/values.local.yaml


gzip -f index.yaml  
gunzip index.yaml.gz