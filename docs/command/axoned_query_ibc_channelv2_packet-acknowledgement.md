## axoned query ibc channelv2 packet-acknowledgement

Query a channel/v2 packet acknowledgement

### Synopsis

Query a channel/v2 packet acknowledgement by client-id and sequence

```
axoned query ibc channelv2 packet-acknowledgement [client-id] [sequence] [flags]
```

### Examples

```
axoned query ibc channelv2 packet-acknowledgement [client-id] [sequence]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for packet-acknowledgement
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
      --prove              show proofs for the query results (default true)
```

### SEE ALSO

* [axoned query ibc channelv2](axoned_query_ibc_channelv2.md)	 - IBC channel/v2 query subcommands
