FROM golang:1.23-alpine AS modules
WORKDIR /workdir/app
COPY app/go.mod app/go.sum ./
RUN go mod download

FROM golang:1.23-alpine AS tester
WORKDIR /workdir/app
COPY ./app ./
RUN go test -v ./...

FROM tester AS builder
RUN go build -o main .

FROM alpine:latest  
WORKDIR /root
COPY --from=builder /workdir/app/main .
EXPOSE 8080
CMD ["./main"]
