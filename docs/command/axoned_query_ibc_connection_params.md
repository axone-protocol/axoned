## axoned query ibc connection params

Query the current ibc connection parameters

### Synopsis

Query the current ibc connection parameters

```
axoned query ibc connection params [flags]
```

### Examples

```
axoned query ibc connection params
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for params
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query ibc connection](axoned_query_ibc_connection.md)	 - IBC connection query subcommands
