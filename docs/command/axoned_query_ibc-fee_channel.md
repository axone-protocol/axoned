## axoned query ibc-fee channel

Query the ibc-fee enabled status of a channel

### Synopsis

Query the ibc-fee enabled status of a channel

```
axoned query ibc-fee channel [port-id] [channel-id] [flags]
```

### Examples

```
axoned query ibc-fee channel transfer channel-6
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for channel
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query ibc-fee](axoned_query_ibc-fee.md)	 - IBC relayer incentivization query subcommands
