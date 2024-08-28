#--- Build stage
FROM golang:1.21-alpine3.18 AS go-builder

WORKDIR /src

# CosmWasm: see https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v2.1.2/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v2.1.2/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a

# hadolint ignore=DL4006
RUN set -eux \
    && apk add --no-cache ca-certificates=20240226-r0 build-base=0.5-r3 git=2.40.1-r0 linux-headers=6.3-r0 \
    && sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 0881c5b463e89e229b06370e9e2961aec0a5c636772d5142c68d351564464a66 \
    && sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 58e1f6bfa89ee390cb9abc69a5bc126029a497fe09dd399f38a82d0d86fe95ef

COPY . /src/

RUN BUILD_TAGS=muslc LINK_STATICALLY=true make build

#--- Image stage
FROM alpine:3.18.6

COPY --from=go-builder /src/target/dist/axoned /usr/bin/axoned

WORKDIR /opt

ENTRYPOINT ["axoned"]
