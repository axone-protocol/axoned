## axoned query ibc client states

Query all available light clients

### Synopsis

Query all available light clients

```
axoned query ibc client states [flags]
```

### Examples

```
axoned query ibc client states
```

### Options

```
      --count-total        count total number of records in client states to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for states
      --limit uint         pagination limit of client states to query for (default 100)
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
      --offset uint        pagination offset of client states to query for
  -o, --output string      Output format (text|json) (default "text")
      --page uint          pagination page of client states to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of client states to query for
      --reverse            results are sorted in descending order
```

### SEE ALSO

* [axoned query ibc client](axoned_query_ibc_client.md)	 - IBC client query subcommands
