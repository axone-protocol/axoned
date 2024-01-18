## okp4d query ibc-fee packets

Query for all of the unrelayed incentivized packets and associated fees across all channels.

### Synopsis

Query for all of the unrelayed incentivized packets and associated fees across all channels.

```
okp4d query ibc-fee packets [flags]
```

### Examples

```
okp4d query ibc-fee packets
```

### Options

```
      --count-total        count total number of records in packets to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for packets
      --limit uint         pagination limit of packets to query for (default 100)
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
      --offset uint        pagination offset of packets to query for
  -o, --output string      Output format (text|json) (default "text")
      --page uint          pagination page of packets to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of packets to query for
      --reverse            results are sorted in descending order
```

### SEE ALSO

* [okp4d query ibc-fee](okp4d_query_ibc-fee.md)	 - IBC relayer incentivization query subcommands
