## okp4d query wasm list-contract-by-code

List wasm all bytecode on the chain for given code id

### Synopsis

List wasm all bytecode on the chain for given code id

```
okp4d query wasm list-contract-by-code [code_id] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-contract-by-code
      --limit uint         pagination limit of list contracts by code to query for (default 100)
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
      --page-key string    pagination page-key of list contracts by code to query for
      --reverse            results are sorted in descending order
```

### SEE ALSO

* [okp4d query wasm](okp4d_query_wasm.md)	 - Querying commands for the wasm module
