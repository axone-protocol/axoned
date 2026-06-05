#!/usr/bin/env bash

set -eo pipefail

# Resolve Go binary path: prefer GOBIN, then GOPATH/bin, finally fall back to `go env GOPATH`/bin.
GOBIN_PATH="${GOBIN:-${GOPATH:+${GOPATH}/bin}}"
GOBIN_PATH="${GOBIN_PATH:-$(go env GOPATH)/bin}"
export PATH="${GOBIN_PATH}:${PATH}"

protoc_install_proto_gen_doc() {
  echo "Installing protobuf protoc-gen-doc plugin"
  (go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest 2> /dev/null)
}

protoc_install_proto_gen_doc

echo "Generating proto docs"
cd proto
for MODULE in $(find . -name '*.proto' -maxdepth 3 -print0 | xargs -0 -n1 dirname | sort | uniq | xargs -n1 dirname); do
  echo "Generating docs for ${MODULE}"
  buf generate --path "${MODULE}" --template buf.gen.doc.yml -v
  mv -f ../docs/proto/docs.md ../docs/proto/"${MODULE}".md
done
