#--- Build stage
FROM golang:1.21-alpine3.17 AS go-builder

WORKDIR /src

# CosmWasm: see https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.5.0/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.5.0/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 2687afbdae1bc6c7c8b05ae20dfb8ffc7ddc5b4e056697d0f37853dfe294e913
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 465e3a088e96fd009a11bfd234c69fb8a0556967677e54511c084f815cf9ce63

# hadolint ignore=DL4006
RUN set -eux \
    && apk add --no-cache ca-certificates=20230506-r0 build-base=0.5-r3 git=2.38.5-r0 linux-headers=5.19.5-r0 \
    && sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 2687afbdae1bc6c7c8b05ae20dfb8ffc7ddc5b4e056697d0f37853dfe294e913 \
    && sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 465e3a088e96fd009a11bfd234c69fb8a0556967677e54511c084f815cf9ce63 \
    && cp "/lib/libwasmvm_muslc.$(uname -m).a" /lib/libwasmvm_muslc.a

COPY . /src/

RUN BUILD_TAGS=muslc LINK_STATICALLY=true make build

#--- Image stage
FROM alpine:3.17.3

COPY --from=go-builder /src/target/dist/okp4d /usr/bin/okp4d

WORKDIR /opt

ENTRYPOINT ["okp4d"]
