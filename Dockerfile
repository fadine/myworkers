FROM golang:1.10.0-alpine3.7 AS builder
ADD . /go/src/github.com/fadine/myworkers
WORKDIR /go/src/github.com/fadine/myworkers
RUN apk update && \
    apk add -U build-base git curl libstdc++ ca-certificates && \
    go env && go list all | grep cover && \
    GOPATH=/go make docker

FROM alpine:3.7
RUN apk add --no-cache --virtual .run-deps \
	curl

COPY ./conf /conf
COPY --from=builder /go/src/github.com/fadine/myworkers/myworkers /myworkers
COPY --from=builder /usr/bin/curl /bin/curl
ENTRYPOINT ["/myworkers"]