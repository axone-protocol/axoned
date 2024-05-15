## axoned query staking historical-info

Query historical info at given height

```
axoned query staking historical-info [height] [flags]
```

### Examples

```
$ axoned query staking historical-info 5
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for historical-info
      --no-indent          Do not indent JSON output
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query staking](axoned_query_staking.md)	 - Querying commands for the staking module
