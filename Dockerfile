# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18-buster AS build

WORKDIR /app

COPY . /app/

RUN CGO_ENABLED=0 GOOS=linux go build -o /mssql_exporter ./cmd

##
## Deploy
##
FROM alpine:3

WORKDIR /

COPY --from=build /mssql_exporter /mssql_exporter
COPY entrypoint.sh /entrypoint.sh

#RUN apk add --no-cache \
#        musl
#
ENTRYPOINT ["/entrypoint.sh"]

