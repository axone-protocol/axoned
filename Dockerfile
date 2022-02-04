#--- Build stage
FROM golang:1.17.6-stretch AS go-builder

WORKDIR /src

COPY . /src/

RUN make build CGO_ENABLED=0 GOOS=linux GOARCH=amd64

#--- Image stage
FROM alpine:3.15

COPY --from=go-builder /src/target/okp4d /usr/bin/okp4d

WORKDIR /opt

CMD ["/usr/bin/okp4d", "version"]
