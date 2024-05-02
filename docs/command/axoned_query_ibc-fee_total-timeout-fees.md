## axoned query ibc-fee total-timeout-fees

Query the total timeout fees for a packet

### Synopsis

Query the total timeout fees for a packet

```
axoned query ibc-fee total-timeout-fees [port-id] [channel-id] [sequence] [flags]
```

### Examples

```
axoned query ibc-fee total-timeout-fees transfer channel-5 100
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for total-timeout-fees
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query ibc-fee](axoned_query_ibc-fee.md)	 - IBC relayer incentivization query subcommands
