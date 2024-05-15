## axoned query bank denom-metadata-by-query-string

Execute the DenomMetadataByQueryString RPC method

```
axoned query bank denom-metadata-by-query-string [flags]
```

### Options

```
      --denom string       
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for denom-metadata-by-query-string
      --no-indent          Do not indent JSON output
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query bank](axoned_query_bank.md)	 - Querying commands for the bank module
