## okp4d query authz grants-by-grantee

query authorization grants granted to a grantee

### Synopsis

Query authorization grants granted to a grantee.
Examples:
$ okp4d q authz grants-by-grantee cosmos1skj..

```
okp4d query authz grants-by-grantee [grantee-addr] [flags]
```

### Options

```
      --count-total        count total number of records in grantee-grants to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for grants-by-grantee
      --limit uint         pagination limit of grantee-grants to query for (default 100)
      --node string        &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --offset uint        pagination offset of grantee-grants to query for
  -o, --output string      Output format (text|json) (default "text")
      --page uint          pagination page of grantee-grants to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of grantee-grants to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query authz](okp4d_query_authz.md)	 - Querying commands for the authz module
