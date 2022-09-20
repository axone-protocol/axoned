FROM bufbuild/buf as BUF

FROM golang:1.19-alpine AS GOLANG
ENV GO111MODULE=on

RUN apk add --no-cache \
  git \
  make

ENV GOLANG_PROTOBUF_VERSION=1.28.0 \
  GRPC_GATEWAY_VERSION=1.16.0

RUN go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v${GOLANG_PROTOBUF_VERSION}
RUN go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v${GRPC_GATEWAY_VERSION}

# install all gogo protobuf binaries
RUN git clone https://github.com/regen-network/protobuf.git; \
    cd protobuf; \
    go mod download; \
    make install

# we need to use git clone because we use 'replace' directive in go.mod
# protoc-gen-gocosmos was moved to to in cosmos/gogoproto but pending a migration there.
RUN git clone https://github.com/regen-network/cosmos-proto.git; \
    cd cosmos-proto/protoc-gen-gocosmos; \
    go install .


COPY --from=BUF /usr/local/bin /usr/local/bin

#FROM ubuntu:14.04
#SHELL ["/bin/sh", "-c"]
#CMD ["bin/sh"]
#COPY --from=golang /go/bin/protoc-gen-go /go/bin/protoc-gen-go-grpc
#COPY --from=buf /usr/local/bin/buf /usr/local/bin/buf
#ENV PATH="/go/bin:${PATH}"
#ENTRYPOINT ["/usr/local/bin/buf"]
