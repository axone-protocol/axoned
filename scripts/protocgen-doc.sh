#!/usr/bin/env bash

set -eo pipefail

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
