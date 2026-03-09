#--- Build stage
FROM golang:1.26-alpine3.22@sha256:169d3991a4f795124a88b33c73549955a3d856e26e8504b5530c30bd245f9f1b AS go-builder

WORKDIR /src

# CosmWasm: see https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v3.0.2/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v3.0.2/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a

SHELL ["/bin/ash", "-o", "pipefail", "-c"]
# hadolint ignore=DL3018
RUN \
    apk add --no-cache ca-certificates build-base=0.5-r3 git=~2.49 linux-headers=6.14.2-r0 \
 && sha256sum /lib/libwasmvm_muslc.aarch64.a | grep b9df5056ab9f61d3f9b944060b44e893d7ade7dad6ff134b36276be0f9a4185a \
 && sha256sum /lib/libwasmvm_muslc.x86_64.a | grep b249396cf884b207f49f46bcf5b8d1fd73b8618eebbe35afb8bf60a8bb24f30a

COPY . /src/

RUN BUILD_TAGS=muslc LINK_STATICALLY=true make build-go

#--- Image stage
FROM alpine:3.22.2

COPY --from=go-builder /src/target/dist/axoned /usr/bin/axoned

WORKDIR /opt

ENTRYPOINT ["axoned"]
