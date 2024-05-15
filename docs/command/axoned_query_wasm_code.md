## axoned query wasm code

Downloads wasm bytecode for given code id

### Synopsis

Downloads wasm bytecode for given code id

```
axoned query wasm code [code_id] [output filename] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for code
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query wasm](axoned_query_wasm.md)	 - Querying commands for the wasm module
