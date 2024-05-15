## axoned query feegrant grant

Query details of a single grant

### Synopsis

Query details for a grant. You can find the fee-grant of a granter and grantee.

```
axoned query feegrant grant [granter] [grantee] [flags]
```

### Examples

```
$ axoned query feegrant grant [granter] [grantee]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for grant
      --no-indent          Do not indent JSON output
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query feegrant](axoned_query_feegrant.md)	 - Querying commands for the feegrant module
