#--- Build stage
FROM golang:1.19-alpine3.16 AS go-builder

WORKDIR /src

# CosmWasm: see https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.2.4/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.2.4/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a

# hadolint ignore=DL4006
RUN set -eux \
    && apk add --no-cache ca-certificates=20220614-r0 build-base=0.5-r3 git=2.36.6-r0 linux-headers=5.16.7-r1 \
    && sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 682a54082e131eaff9beec80ba3e5908113916fcb8ddf7c668cb2d97cb94c13c \
    && sha256sum /lib/libwasmvm_muslc.x86_64.a | grep ce3d892377d2523cf563e01120cb1436f9343f80be952c93f66aa94f5737b661 \
    && cp "/lib/libwasmvm_muslc.$(uname -m).a" /lib/libwasmvm_muslc.a

COPY . /src/

RUN BUILD_TAGS=muslc LINK_STATICALLY=true make build

#--- Image stage
FROM alpine:3.17.3

COPY --from=go-builder /src/target/dist/okp4d /usr/bin/okp4d

WORKDIR /opt

ENTRYPOINT ["okp4d"]
