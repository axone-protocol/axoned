## axoned query ibc-transfer total-escrow

Query the total amount of tokens in escrow for a denom

### Synopsis

Query the total amount of tokens in escrow for a denom

```
axoned query ibc-transfer total-escrow [denom] [flags]
```

### Examples

```
axoned query ibc-transfer total-escrow uosmo
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for total-escrow
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query ibc-transfer](axoned_query_ibc-transfer.md)	 - IBC fungible token transfer query subcommands
