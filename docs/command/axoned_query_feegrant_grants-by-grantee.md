## axoned query feegrant grants-by-grantee

Query all grants of a grantee

### Synopsis

Queries all the grants for a grantee address.

```
axoned query feegrant grants-by-grantee [grantee] [flags]
```

### Examples

```
$ axoned query feegrant grants-by-grantee [grantee]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for grants-by-grantee
      --no-indent          Do not indent JSON output
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
      --page-count-total   
      --page-key binary    
      --page-limit uint    
      --page-offset uint   
      --page-reverse       
```

### SEE ALSO

* [axoned query feegrant](axoned_query_feegrant.md)	 - Querying commands for the feegrant module
