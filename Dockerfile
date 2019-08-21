FROM golang:1.12-alpine
RUN apk update
RUN apk add openssl ca-certificates git curl
ENV GO111MODULE on
WORKDIR /go/src/github.com/akkeris/service-watcher-f5
COPY . .
RUN mkdir -p /root/.kube/certs
CMD "/start.sh"



