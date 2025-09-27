## axoned query ibc channelv2 unreceived-packets

Query a channel/v2 unreceived-packets

### Synopsis

Query a channel/v2 unreceived-packets by client-id and sequences

```
axoned query ibc channelv2 unreceived-packets [client-id] [flags]
```

### Examples

```
axoned query ibc channelv2 unreceived-packet [client-id] --sequences=1,2,3
```

### Options

```
      --grpc-addr string       the gRPC endpoint to use for this chain
      --grpc-insecure          allow gRPC over insecure channels, if not the server must use TLS
      --height int             Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                   help for unreceived-packets
      --node string            <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string          Output format (text|json) (default "text")
      --sequences int64Slice   comma separated list of packet sequence numbers (default [])
```

### SEE ALSO

* [axoned query ibc channelv2](axoned_query_ibc_channelv2.md)	 - IBC channel/v2 query subcommands
