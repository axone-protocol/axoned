## okp4d query wasm list-contracts-by-creator

List all contracts by creator

### Synopsis

List all contracts by creator

```
okp4d query wasm list-contracts-by-creator [creator] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-contracts-by-creator
      --limit uint         pagination limit of list contracts by creator to query for (default 100)
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
      --page-key string    pagination page-key of list contracts by creator to query for
      --reverse            results are sorted in descending order
```

### SEE ALSO

* [okp4d query wasm](okp4d_query_wasm.md)	 - Querying commands for the wasm module
