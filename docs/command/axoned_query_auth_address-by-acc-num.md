## axoned query auth address-by-acc-num

Query account address by account number

```
axoned query auth address-by-acc-num [acc-num] [flags]
```

### Options

```
      --account-id uint    
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for address-by-acc-num
      --no-indent          Do not indent JSON output
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query auth](axoned_query_auth.md)	 - Querying commands for the auth module
