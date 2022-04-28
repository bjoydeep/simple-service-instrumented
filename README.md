# Simple Service

### Service content
a-service is a simple golang service which exposes 2 http endpoints:
1. /health
1. /metrics

It also exposes a bunch of metrics to prometheus:
1. regular golang metrics
1. a_hits_total: The total number of hits
1. avec_duration_seconds: 99th percentile latency in seconds

You can run a docker build and post it to your own docker registry. However for the sake of simplicity a docker images is already available with this a-service at quay.io/bjoydeep/a-service:latest

### Goal
Goal is to get the metrics this service emits into the platform Prometheus inside OpenShift - in other words, the Promtheus deployed by the default OpenShift. _**Check with OpenShift RedHat Team if this will be supported even though you can technically do it.**_

#### Steps
1. Connect to your OpenShift cluster as cluster admin.
1. Adjust image location in deploy/deploy.yaml if needed.
1. Run `kubectl apply -k deploy/`
1. You will find that a namespace a-service created and a pod and a service is created under it. It will also create `ServiceMonitor CR` in the `openshift-monitoring` namespace and allow `prometheus-k8s` service account in openshift-monitoring namespace to scrape endpoints of service exposed by `a-service` (it creates the required ClusterRole and ClusterRoleBinding)
1. Create a route in OpenShift for the service created by a-service
1. Open up a browser and make requests to http://a-service.*/health a few times.
1. Check local prometheus for metrics mentioned above