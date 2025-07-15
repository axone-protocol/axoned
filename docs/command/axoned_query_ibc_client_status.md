## axoned query ibc client status

Query client status

### Synopsis

Query client activity status. Any client without an 'Active' status is considered inactive

```
axoned query ibc client status [client-id] [flags]
```

### Examples

```
axoned query ibc client status [client-id]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for status
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query ibc client](axoned_query_ibc_client.md)	 - IBC client query subcommands
