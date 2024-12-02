FROM golang:1.17-alpine

RUN apk update && apk add --no-cache iptables bash

WORKDIR /app

COPY . .

RUN go build -o trustify main.go

ENTRYPOINT ["/app/trustify"]