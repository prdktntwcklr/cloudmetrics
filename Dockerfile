FROM golang:1.23-alpine AS builder

WORKDIR /workdir/app
COPY app/ .
RUN go build -o main .

FROM alpine:latest  
WORKDIR /root

COPY --from=builder /workdir/app/main .

EXPOSE 8080
CMD ["./main"]
