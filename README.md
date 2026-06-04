# Cloudmetrics

A Go-based sensor simulator for Kubernetes, featuring a complete observability
stack. It generates real-time temperature telemetry and reading counts, exported
via Prometheus metrics and Loki-compatible logs, fully automated and deployed
via GitOps using Argo CD.

## Prerequisites

- A running [Kubernetes](https://kubernetes.io/) cluster
- [Helm](https://helm.sh/) (package manager for Kubernetes) installed

## Deployment

### Install Argo CD

First, install Argo CD into your cluster:

```bash
kubectl create namespace argocd
kubectl apply -n argocd --server-side --force-conflicts -f https://raw.githubusercontent.com/argoproj/argo-cd/v3.4.3/manifests/install.yaml
```

### Create the Grafana Password 

Before deploying, create the password required to log into the Grafana UI.
Ensure that this namespace matches the target namespace of your Grafana
deployment (defaulting to `monitoring`).

```bash
kubectl create namespace monitoring
kubectl create secret generic grafana-admin-secret \
  -n monitoring \
  --from-literal=admin-user=admin \
  --from-literal=admin-password='supersecurepassword'
```

### Deploy the Stack (App-of-Apps)

This uses the Argo CD **App of Apps** pattern. Deploying the single root
bootstrap application will automatically discover, configure, and roll out the
entire application and observability pipeline:

```bash
kubectl apply -f argocd-apps/root-app.yaml
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

    **Logs**: Go to `Explore`, select the `Loki` data source, and run a query
    such as: `{service_name="cloudmetrics-app"}`.

    ![Grafana Logs](images/logs.png)

3. **Argo CD Dashboard**: Use port forwarding:

    ```bash
    # Use local port 7070 to avoid conflicts with the application port (8080)
    kubectl port-forward svc/argocd-server -n argocd 7070:443
    ```

    After that, open your browser and head to: https://localhost:7070/. Ignore
    the Certificate Error warning (Argo CD uses a self-signed certificate) and
    click on `Proceed anyway`. 

    The default username is `admin`. Argo CD automatically generates a unique
    startup password and stores it securely in a Kubernetes secret. Retrieve and
    decode it by running:

    ```bash
    kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo
    ```

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

To remove the observability stack:

```bash
helmfile -f deploy/observability/helmfile.yaml destroy
kubectl delete ns monitoring
```

Note that Custom Resource Definitions (CRDs) are not deleted automatically. For
details on how to uninstall them, refer to the [kube-prometheus-stack official documentation](https://github.com/prometheus-community/helm-charts/tree/556bfa39ea386b9d261b5ca49a9dc62f112ec78f/charts/kube-prometheus-stack#uninstall-helm-chart).
