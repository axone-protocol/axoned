## axoned query ibc channel client-state

Query the client state associated with a channel

### Synopsis

Query the client state associated with a channel, by providing its port and channel identifiers.

```
axoned query ibc channel client-state [port-id] [channel-id] [flags]
```

### Examples

```
axoned query ibc channel client-state [port-id] [channel-id]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for client-state
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query ibc channel](axoned_query_ibc_channel.md)	 - IBC channel query subcommands
