## axoned query interchain-accounts host packet-events

Query the interchain-accounts host submodule packet events

### Synopsis

Query the interchain-accounts host submodule packet events for a particular channel and sequence

```
axoned query interchain-accounts host packet-events [channel-id] [sequence] [flags]
```

### Examples

```
axoned query interchain-accounts host packet-events channel-0 100
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for packet-events
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query interchain-accounts host](axoned_query_interchain-accounts_host.md)	 - IBC interchain accounts host query subcommands
