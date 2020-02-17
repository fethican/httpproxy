FROM golang:1.13.8-alpine3.11 AS builder
RUN apk --no-cache add gcc musl-dev git
RUN go get -u github.com/fethican/httpproxy
WORKDIR /go/src/github.com/fethican/httpproxy
ENV APP_VERSION=v0.1
RUN git checkout ${APP_VERSION} > /dev/null 2>&1
RUN githash=$(git rev-parse --short HEAD 2>/dev/null) \
    && today=$(date +%Y-%m-%d --utc) \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags '-s -w -X main.Version=${APP_VERSION} -X main.Commit=${githash} -X main.BuiltAt=${today}' \
    -o /binary

FROM alpine:3.11 AS libs
RUN apk --no-cache add ca-certificates

FROM scratch
COPY --from=libs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /binary /httpproxy
ENTRYPOINT ["/httpproxy"]
