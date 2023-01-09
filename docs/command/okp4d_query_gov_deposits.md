## okp4d query gov deposits

Query deposits on a proposal

### Synopsis

Query details for all deposits on a proposal.
You can find the proposal-id by running "okp4d query gov proposals".

Example:
$ okp4d query gov deposits 1

```
okp4d query gov deposits [proposal-id] [flags]
```

### Options

```
      --count-total        count total number of records in deposits to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for deposits
      --limit uint         pagination limit of deposits to query for (default 100)
      --node string        &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --offset uint        pagination offset of deposits to query for
  -o, --output string      Output format (text|json) (default "text")
      --page uint          pagination page of deposits to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of deposits to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query gov](okp4d_query_gov.md)	 - Querying commands for the governance module
