#--- Build stage
FROM golang:1.27.0-alpine3.23@sha256:3747dcba41c8b0db3211fda4db61638b980e17ac5bb3c94460a975a9cfe19395 AS go-builder

WORKDIR /src

# CosmWasm: see https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v3.0.8/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v3.0.8/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a

SHELL ["/bin/ash", "-o", "pipefail", "-c"]
# hadolint ignore=DL3018
RUN \
    apk add --no-cache ca-certificates=20260909-r0 build-base=0.5-r3 git=2.52.0-r0 linux-headers=6.16.12-r0 \
 && sha256sum /lib/libwasmvm_muslc.aarch64.a | grep c73a0d5d340e35188e138584ddd6662a160902adef1b08b209e38d16b43a4c28 \
 && sha256sum /lib/libwasmvm_muslc.x86_64.a | grep b2299c85d49faccf3dcbb84984f30f55e8870111df98c10f017f86204d007470 \
 && install -d -m 1777 /empty-tmp

COPY . /src/

RUN BUILD_TAGS=muslc LINK_STATICALLY=true make build-go

#--- Image stage
FROM scratch

COPY --from=go-builder /src/target/dist/axoned /usr/bin/axoned
COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=go-builder /empty-tmp /tmp

WORKDIR /opt

ENTRYPOINT ["axoned"]
