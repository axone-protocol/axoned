## axoned query ibc channel unreceived-packets

Query all the unreceived packets associated with a channel

### Synopsis

Determine if a packet, given a list of packet commitment sequences, is unreceived.

The return value represents:

- Unreceived packet commitments: no acknowledgement exists on receiving chain for the given packet commitment sequence on sending chain.

```
axoned query ibc channel unreceived-packets [port-id] [channel-id] [flags]
```

### Examples

```
axoned query ibc channel unreceived-packets [port-id] [channel-id] --sequences=1,2,3
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

- [axoned query ibc channel](axoned_query_ibc_channel.md)	 - IBC channel query subcommands
