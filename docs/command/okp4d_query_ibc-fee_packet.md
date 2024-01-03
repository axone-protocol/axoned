## okp4d query ibc-fee packet

Query for an unrelayed incentivized packet by port-id, channel-id and packet sequence.

### Synopsis

Query for an unrelayed incentivized packet by port-id, channel-id and packet sequence.

```
okp4d query ibc-fee packet [port-id] [channel-id] [sequence] [flags]
```

### Examples

```
okp4d query ibc-fee packet
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for packet
      --node string        <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query ibc-fee](okp4d_query_ibc-fee.md)	 - IBC relayer incentivization query subcommands
