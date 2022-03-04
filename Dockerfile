#--- Build stage
FROM golang:1.18rc1-stretch AS go-builder

WORKDIR /src

COPY . /src/

RUN make build CGO_ENABLED=0

#--- Image stage
FROM alpine:3.15

COPY --from=go-builder /src/target/dist/okp4d /usr/bin/okp4d

WORKDIR /opt

CMD ["/usr/bin/okp4d", "version"]
