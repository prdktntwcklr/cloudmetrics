# Cloudmetrics

A Go-based sensor simulator designed to run entirely in Docker.

## Quick Start

Run the following command in the root of the project to generate the `go.mod` and
`go.sum` files:

```bash
docker run --rm -v ./app:/app -w /app golang:1.23-alpine sh -c "go mod init cloudmetrics-app && go mod tidy"
```

Build the optimized container using the multi-stage Dockerfile, which also
includes unit tests:

```bash
docker build -t cloudmetrics-app .
```

Start the container and map the port to your local machine:

```bash
docker run --rm -p 8080:8080 cloudmetrics-app
```

Navigate to http://localhost:8080/metrics in your browser to observe the live
sensor data.
