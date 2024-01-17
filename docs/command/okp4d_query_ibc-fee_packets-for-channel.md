## okp4d query ibc-fee packets-for-channel

Query for all of the unrelayed incentivized packets on a given channel

### Synopsis

Query for all of the unrelayed incentivized packets on a given channel. These are packets that have not yet been relayed.

```
okp4d query ibc-fee packets-for-channel [port-id] [channel-id] [flags]
```

### Examples

```
okp4d query ibc-fee packets-for-channel
```

### Options

```
      --count-total        count total number of records in packets-for-channel to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for packets-for-channel
      --limit uint         pagination limit of packets-for-channel to query for (default 100)
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
      --offset uint        pagination offset of packets-for-channel to query for
  -o, --output string      Output format (text|json) (default "text")
      --page uint          pagination page of packets-for-channel to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of packets-for-channel to query for
      --reverse            results are sorted in descending order
```

### SEE ALSO

* [okp4d query ibc-fee](okp4d_query_ibc-fee.md)	 - IBC relayer incentivization query subcommands
