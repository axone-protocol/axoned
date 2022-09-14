#--- Build stage
FROM golang:1.19-alpine3.16 AS go-builder

WORKDIR /src

# hadolint ignore=DL4006
RUN set -eux \
    && apk add --no-cache ca-certificates=20220614-r0 build-base=0.5-r3 git=2.36.2-r0

COPY . /src/

RUN BUILD_TAGS=muslc LINK_STATICALLY=true make build

#--- Image stage
FROM alpine:3.16.2

COPY --from=go-builder /src/target/dist/okp4d /usr/bin/okp4d

WORKDIR /opt

ENTRYPOINT ["okp4d"]
