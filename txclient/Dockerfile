FROM golang:1.20-alpine3.16 AS builder

WORKDIR /usr/local/go/src/

ADD .. /usr/local/go/src/

RUN go mod download
RUN go build -mod=readonly -o client cmd/main.go

FROM alpine:latest

COPY --from=builder /usr/local/go/src/client /

EXPOSE ${APP_PORT}
CMD ["/client"]