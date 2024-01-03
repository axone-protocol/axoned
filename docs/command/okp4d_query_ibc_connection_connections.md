## okp4d query ibc connection connections

Query all connections

### Synopsis

Query all connections ends from a chain

```
okp4d query ibc connection connections [flags]
```

### Examples

```
okp4d query ibc connection connections
```

### Options

```
      --count-total        count total number of records in connection ends to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for connections
      --limit uint         pagination limit of connection ends to query for (default 100)
      --node string        <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --offset uint        pagination offset of connection ends to query for
  -o, --output string      Output format (text|json) (default "text")
      --page uint          pagination page of connection ends to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of connection ends to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query ibc connection](okp4d_query_ibc_connection.md)	 - IBC connection query subcommands
