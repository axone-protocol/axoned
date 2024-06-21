#--- Build stage
FROM golang:1.21-alpine3.18 AS go-builder

WORKDIR /src

# CosmWasm: see https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v2.0.0/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v2.0.0/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a

# hadolint ignore=DL4006
RUN set -eux \
    && apk add --no-cache ca-certificates=20240226-r0 build-base=0.5-r3 git=2.40.1-r0 linux-headers=6.3-r0 \
    && sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 3b478b3e51d31e53ce9324a8895d2cd7278af5179b9a02ea55d8627958e42afa \
    && sha256sum /lib/libwasmvm_muslc.x86_64.a | grep ca08bb7b73b49b483611d9755bb8455620bb8c0faf3014400908ed49bf3b19a5

COPY . /src/

RUN BUILD_TAGS=muslc LINK_STATICALLY=true make build

#--- Image stage
FROM alpine:3.20.1

COPY --from=go-builder /src/target/dist/axoned /usr/bin/axoned

WORKDIR /opt

ENTRYPOINT ["axoned"]
