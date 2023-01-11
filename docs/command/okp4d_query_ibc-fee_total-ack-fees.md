## okp4d query ibc-fee total-ack-fees

Query the total acknowledgement fees for a packet

### Synopsis

Query the total acknowledgement fees for a packet

```
okp4d query ibc-fee total-ack-fees [port-id] [channel-id] [sequence] [flags]
```

### Examples

```
okp4d query ibc-fee total-ack-fees transfer channel-5 100
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for total-ack-fees
      --node string        &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query ibc-fee](okp4d_query_ibc-fee.md)	 - IBC relayer incentivization query subcommands
