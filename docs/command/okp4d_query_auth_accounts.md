## okp4d query auth accounts

Query all the accounts

```
okp4d query auth accounts [flags]
```

### Options

```
      --count-total        count total number of records in all-accounts to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for accounts
      --limit uint         pagination limit of all-accounts to query for (default 100)
      --node string        &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --offset uint        pagination offset of all-accounts to query for
  -o, --output string      Output format (text|json) (default "text")
      --page uint          pagination page of all-accounts to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of all-accounts to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query auth](okp4d_query_auth.md)	 - Querying commands for the auth module
