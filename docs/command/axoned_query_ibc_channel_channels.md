## axoned query ibc channel channels

Query all channels

### Synopsis

Query all channels from a chain

```
axoned query ibc channel channels [flags]
```

### Examples

```
axoned query ibc channel channels
```

### Options

```
      --count-total        count total number of records in channels to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for channels
      --limit uint         pagination limit of channels to query for (default 100)
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
      --offset uint        pagination offset of channels to query for
  -o, --output string      Output format (text|json) (default "text")
      --page uint          pagination page of channels to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of channels to query for
      --reverse            results are sorted in descending order
```

### SEE ALSO

* [axoned query ibc channel](axoned_query_ibc_channel.md)	 - IBC channel query subcommands
