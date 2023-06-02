#--- Build stage
FROM golang:1.19-alpine3.16 AS go-builder

WORKDIR /src

# CosmWasm: see https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.2.3/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.2.3/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a

# hadolint ignore=DL4006
RUN set -eux \
    && apk add --no-cache ca-certificates=20220614-r0 build-base=0.5-r3 git=2.36.6-r0 linux-headers=5.16.7-r1 \
    && sha256sum /lib/libwasmvm_muslc.aarch64.a | grep d6904bc0082d6510f1e032fc1fd55ffadc9378d963e199afe0f93dd2667c0160 \
    && sha256sum /lib/libwasmvm_muslc.x86_64.a | grep bb8ffda690b15765c396266721e45516cb3021146fd4de46f7daeda5b0d82c86 \
    && cp "/lib/libwasmvm_muslc.$(uname -m).a" /lib/libwasmvm_muslc.a

COPY . /src/

RUN BUILD_TAGS=muslc LINK_STATICALLY=true make build

#--- Image stage
FROM alpine:3.17.3

COPY --from=go-builder /src/target/dist/okp4d /usr/bin/okp4d

WORKDIR /opt

ENTRYPOINT ["okp4d"]
