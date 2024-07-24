#--- Build stage
FROM golang:1.21-alpine3.18 AS go-builder

WORKDIR /src

# CosmWasm: see https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v2.1.0/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v2.1.0/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a

# hadolint ignore=DL4006
RUN set -eux \
    && apk add --no-cache ca-certificates=20240226-r0 build-base=0.5-r3 git=2.40.1-r0 linux-headers=6.3-r0 \
    && sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 06945cae8fced839a2be0980887a1c5d04d15fd29837ac644a648d555c49ab4d \
    && sha256sum /lib/libwasmvm_muslc.x86_64.a | grep df4bd912c35be48781a40edea88fd5f409c643fb27e0dc043184ef51dc50a1cc

COPY . /src/

RUN BUILD_TAGS=muslc LINK_STATICALLY=true make build

#--- Image stage
FROM alpine:3.20.2

COPY --from=go-builder /src/target/dist/axoned /usr/bin/axoned

WORKDIR /opt

ENTRYPOINT ["axoned"]
