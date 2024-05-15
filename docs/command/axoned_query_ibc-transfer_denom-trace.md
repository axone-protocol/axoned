## axoned query ibc-transfer denom-trace

Query the denom trace info from a given trace hash or ibc denom

### Synopsis

Query the denom trace info from a given trace hash or ibc denom

```
axoned query ibc-transfer denom-trace [hash/denom] [flags]
```

### Examples

```
axoned query ibc-transfer denom-trace 27A6394C3F9FF9C9DCF5DFFADF9BB5FE9A37C7E92B006199894CF1824DF9AC7C
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for denom-trace
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query ibc-transfer](axoned_query_ibc-transfer.md)	 - IBC fungible token transfer query subcommands
