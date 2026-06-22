# Stage 1: Cache Go modules to speed up subsequent builds.
# This layer is only rebuilt when go.mod or go.sum changes.
FROM golang:1.23-alpine AS modules
WORKDIR /workdir/app
COPY go.mod go.sum ./
RUN go mod download

# Stage 2: Run unit and integration tests
FROM golang:1.23-alpine AS tester
WORKDIR /workdir/app
COPY . .
RUN go test -v ./...

# Stage 3: Build the application
FROM tester AS builder
ARG GIT_SHA
RUN : "${GIT_SHA:?ERROR: GIT_SHA build argument is required!}" && \
    CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.Version=${GIT_SHA}" \
    -o main .

# Stage 4: Package the application into a minimal production image
# TODO: Do not run application as root user
FROM alpine:3.22.4  
WORKDIR /root
COPY --from=builder /workdir/app/main .
EXPOSE 8080
CMD ["./main"]
