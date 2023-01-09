## okp4d query ibc-fee channels

Query the ibc-fee enabled channels

### Synopsis

Query the ibc-fee enabled channels

```
okp4d query ibc-fee channels [flags]
```

### Examples

```
okp4d query ibc-fee channels
```

### Options

```
      --count-total        count total number of records in channels to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for channels
      --limit uint         pagination limit of channels to query for (default 100)
      --node string        &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --offset uint        pagination offset of channels to query for
  -o, --output string      Output format (text|json) (default "text")
      --page uint          pagination page of channels to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of channels to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query ibc-fee](okp4d_query_ibc-fee.md)	 - IBC relayer incentivization query subcommands
