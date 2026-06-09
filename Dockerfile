FROM golang:1.23-alpine AS modules
WORKDIR /workdir/app
COPY app/go.mod app/go.sum ./
RUN go mod download

FROM golang:1.23-alpine AS tester
WORKDIR /workdir/app
COPY ./app ./
RUN go test -v -tags=integration ./...

FROM tester AS builder
ARG GIT_SHA
RUN : "${GIT_SHA:?ERROR: GIT_SHA build argument is required!}"
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.GitCommit=${GIT_SHA}" \
    -o main .

FROM alpine:3.22.4  
WORKDIR /root
COPY --from=builder /workdir/app/main .
COPY ./app/index.html .
EXPOSE 8080
CMD ["./main"]
