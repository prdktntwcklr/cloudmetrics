# Cloudmetrics

A Go-based sensor simulator designed to run in Docker and Kubernetes. It
provides real-time simulated ambient temperature data and reading counts via a 
Prometheus-compatible `/metrics` endpoint.

## Local Build

Build the optimized container using the multi-stage Dockerfile, which also
automatically runs unit tests:

```bash
docker build -t cloudmetrics:latest .
```

## Kubernetes Deployment

To deploy the application to your local Kubernetes cluster (e.g. Docker
Desktop), first deploy the manifests:

```bash
kubectl apply -f k8s/
```

Verify the deployment:

```bash
kubectl get all -l app=cloudmetrics
```

Once the `EXTERNAL-IP` of the service appears, navigate to http://localhost:8080/metrics.

## Development & Testing

If you add new Go packages, update the `go.mod` and `go.sum` files using a
container:

```bash
docker run --rm -v ./app:/app -w /app golang:1.23-alpine sh -c "go mod init cloudmetrics-app ; go mod tidy"
```

## Cleanup

To stop and remove all Kubernetes resources:

```bash
kubectl delete -f k8s/
```
