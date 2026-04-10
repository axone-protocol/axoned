#--- Build stage
FROM golang:1.25-alpine3.23@sha256:7a00384194cf2cb68924bbb918d675f1517357433c8541bac0ab2f929b9d5447 AS go-builder

WORKDIR /src

# CosmWasm: see https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v3.0.2/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v3.0.2/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a

SHELL ["/bin/ash", "-o", "pipefail", "-c"]
# hadolint ignore=DL3018
RUN \
    apk add --no-cache ca-certificates=20251003-r0 build-base=0.5-r3 git=2.52.0-r0 linux-headers=6.16.12-r0 \
 && sha256sum /lib/libwasmvm_muslc.aarch64.a | grep b9df5056ab9f61d3f9b944060b44e893d7ade7dad6ff134b36276be0f9a4185a \
 && sha256sum /lib/libwasmvm_muslc.x86_64.a | grep b249396cf884b207f49f46bcf5b8d1fd73b8618eebbe35afb8bf60a8bb24f30a

COPY . /src/

RUN BUILD_TAGS=muslc LINK_STATICALLY=true make build-go

#--- Image stage
FROM scratch

COPY --from=go-builder /src/target/dist/axoned /usr/bin/axoned
COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /opt

ENTRYPOINT ["axoned"]
