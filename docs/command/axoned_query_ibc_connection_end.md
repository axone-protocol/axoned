## axoned query ibc connection end

Query stored connection end

### Synopsis

Query stored connection end

```
axoned query ibc connection end [connection-id] [flags]
```

### Examples

```
axoned query ibc connection end [connection-id]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for end
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
      --prove              show proofs for the query results (default true)
```

### SEE ALSO

* [axoned query ibc connection](axoned_query_ibc_connection.md)	 - IBC connection query subcommands
