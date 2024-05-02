## axoned query gov vote

Query details of a single vote

```
axoned query gov vote [proposal-id] [voter-addr] [flags]
```

### Examples

```
axoned query gov vote 1 cosmos1...
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for vote
      --no-indent          Do not indent JSON output
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query gov](axoned_query_gov.md)	 - Querying commands for the gov module
