## axoned query wasm contract-history

Prints out the code history for a contract given its address

### Synopsis

Prints out the code history for a contract given its address

```
axoned query wasm contract-history [bech32_address] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for contract-history
      --limit uint         pagination limit of contract history to query for (default 100)
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
      --page-key string    pagination page-key of contract history to query for
      --reverse            results are sorted in descending order
```

### SEE ALSO

* [axoned query wasm](axoned_query_wasm.md)	 - Querying commands for the wasm module
