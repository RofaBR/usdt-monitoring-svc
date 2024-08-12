FROM golang:1.20-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/RofaBR/usdt-monitoring-svc
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/usdt-monitoring-svc /go/src/github.com/RofaBR/usdt-monitoring-svc


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/usdt-monitoring-svc /usr/local/bin/usdt-monitoring-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["usdt-monitoring-svc"]
