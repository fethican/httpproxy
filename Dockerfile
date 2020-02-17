FROM golang:1.13.8-alpine3.11 AS builder
RUN apk --no-cache add gcc musl-dev git make
RUN go get -u github.com/fethican/httpproxy
WORKDIR /go/src/github.com/fethican/httpproxy
ENV APP_VERSION=v0.1
RUN git checkout ${APP_VERSION} > /dev/null 2>&1
RUN make binary

FROM alpine:3.11 AS libs
RUN apk --no-cache add ca-certificates

FROM scratch
COPY --from=libs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /binary /httpproxy
ENTRYPOINT ["/httpproxy"]
