## axoned query ibc client creator

Query a client's creator

### Synopsis

Query a client's creator

```
axoned query ibc client creator [client-id] [flags]
```

### Examples

```
axoned query ibc client creator 08-wasm-0
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for creator
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query ibc client](axoned_query_ibc_client.md)	 - IBC client query subcommands
