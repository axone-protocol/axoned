## axoned query ibc channel upgrade-error

Query the upgrade error

### Synopsis

Query the upgrade error for a given channel

```
axoned query ibc channel upgrade-error [port-id] [channel-id] [flags]
```

### Examples

```
axoned query ibc channel upgrade-error [port-id] [channel-id]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for upgrade-error
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
      --prove              show proofs for the query results (default true)
```

### SEE ALSO

* [axoned query ibc channel](axoned_query_ibc_channel.md)	 - IBC channel query subcommands
