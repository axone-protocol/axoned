## axoned query ibc channelv2 unreceived-acks

Query all the unreceived acks associated with a client

### Synopsis

Given a list of acknowledgement sequences from counterparty, determine if an ack on the counterparty chain has been received on the executing chain.

The return value represents:

- Unreceived packet acknowledgement: packet commitment exists on original sending (executing) chain and ack exists on receiving chain.

```
axoned query ibc channelv2 unreceived-acks [client-id] [flags]
```

### Examples

```
axoned query ibc channelv2 unreceived-acks [client-id] --sequences=1,2,3
```

### Options

```
      --grpc-addr string       the gRPC endpoint to use for this chain
      --grpc-insecure          allow gRPC over insecure channels, if not the server must use TLS
      --height int             Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                   help for unreceived-acks
      --node string            <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string          Output format (text|json) (default "text")
      --sequences int64Slice   comma separated list of packet sequence numbers (default [])
```

### SEE ALSO

- [axoned query ibc channelv2](axoned_query_ibc_channelv2.md)	 - IBC channel/v2 query subcommands
