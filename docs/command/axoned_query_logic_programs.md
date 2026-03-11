## axoned query logic programs

lists stored programs

```
axoned query logic programs [flags]
```

### Examples

```
$ axoned query logic programs --limit 10
```

### Options

```
      --count-total        count total number of records in programs to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for programs
      --limit uint         pagination limit of programs to query for (default 100)
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
      --offset uint        pagination offset of programs to query for
  -o, --output string      Output format (text|json) (default "text")
      --page uint          pagination page of programs to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of programs to query for
      --reverse            results are sorted in descending order
```

### SEE ALSO

* [axoned query logic](axoned_query_logic.md)	 - Querying commands for the logic module
