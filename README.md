# Cloudmetrics

A Go-based sensor simulator designed to run in Docker and Kubernetes. It
provides real-time simulated ambient temperature data and reading counts via a 
Prometheus-compatible `/metrics` endpoint.

## Prerequisites

Before deploying the CloudMetrics app, ensure you have the following:

- A running [Kubernetes](https://kubernetes.io/) cluster (by e.g. using [Docker Desktop](https://www.docker.com/products/docker-desktop/))
- [Helm](https://helm.sh/) (package manager for Kubernetes) installed

The application uses a `ServiceMonitor` to export metrics. This requires the
Prometheus Operator to be installed in your cluster:

```bash
# Add the prometheus-community helm repo
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana-community https://grafana-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

# Create a hardcoded secret to log into Grafana (not for production!)
kubectl create secret generic grafana-admin-secret \
  -n monitoring \
  --from-literal=admin-user=admin \
  --from-literal=admin-password='supersecurepassword'

# Install the kube-prometheus-stack inside the 'monitoring' namespace
# This installs Prometheus, Grafana, and the required Custom Resource Definitions (CRDs)
helm install prometheus prometheus-community/kube-prometheus-stack \
  --create-namespace \
  --namespace monitoring \
  -f kube-prometheus-stack-values.yaml

# Install Loki for log aggregation
helm install loki grafana-community/loki \
  --namespace monitoring \
  -f loki-values.yaml

# Install Alloy to gather and send logs to Loki
helm install alloy grafana/alloy \
  --namespace monitoring \
  -f alloy-values.yaml
```

## Local Build

Build the optimized container using the multi-stage Dockerfile, which also
automatically runs unit tests:

```bash
docker build -t cloudmetrics:latest .
```

## Kubernetes Deployment

To deploy the application to your local Kubernetes cluster, use Helm:

```bash
helm install cloudmetrics-dev charts/cloudmetrics
```

## Accessing Data

1. **Application Metrics**: Once the `EXTERNAL-IP` of the `cloudmetrics-service`
appears, visit http://localhost:8080/metrics to see the raw Prometheus format.

2. **Grafana UI**: To access the Grafana UI, first port-forward it:

    ```bash
    kubectl port-forward svc/prometheus-grafana -n monitoring 3000:80
    ```

    **Metrics**: Open http://localhost:3000/ and log in using the username and
    password you stored in the Kubernetes secret above. You should be able to
    open the `Cloudmetrics` dashboard to see the metrics being displayed.

    ![Grafana UI](images/grafana.png)

    **Logs**: Go to `Explore`, select the `Loki` data source, and run: `{service_name="cloudmetrics-app"}`.

    ![Grafana Logs](images/logs.png)

## Development & Testing

If you add new Go packages, update the `go.mod` and `go.sum` files using a
container:

```bash
docker run --rm -v ./app:/app -w /app golang:1.23-alpine sh -c "go mod init cloudmetrics-app ; go mod tidy"
```

## Cleanup

To stop and remove the application and its associated resources:

```bash
helm uninstall cloudmetrics-dev
```

For details on how to uninstall the `kube-prometheus-stack` chart, refer to the
[official documentation](https://github.com/prometheus-community/helm-charts/tree/556bfa39ea386b9d261b5ca49a9dc62f112ec78f/charts/kube-prometheus-stack#uninstall-helm-chart).
