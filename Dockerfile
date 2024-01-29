#--- Build stage
FROM golang:1.21-alpine3.17 AS go-builder

WORKDIR /src

# CosmWasm: see https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.5.2/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.5.2/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a

# hadolint ignore=DL4006
RUN set -eux \
    && apk add --no-cache ca-certificates=20230506-r0 build-base=0.5-r3 git=2.38.5-r0 linux-headers=5.19.5-r0 \
    && sha256sum /lib/libwasmvm_muslc.aarch64.a | grep e78b224c15964817a3b75a40e59882b4d0e06fd055b39514d61646689cef8c6e \
    && sha256sum /lib/libwasmvm_muslc.x86_64.a | grep e660a38efb2930b34ee6f6b0bb12730adccb040b6ab701b8f82f34453a426ae7 \
    && cp "/lib/libwasmvm_muslc.$(uname -m).a" /lib/libwasmvm_muslc.a

COPY . /src/

RUN BUILD_TAGS=muslc LINK_STATICALLY=true make build

#--- Image stage
FROM alpine:3.19.1

COPY --from=go-builder /src/target/dist/okp4d /usr/bin/okp4d

WORKDIR /opt

ENTRYPOINT ["okp4d"]
