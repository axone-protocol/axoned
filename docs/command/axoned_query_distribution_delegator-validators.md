## axoned query distribution delegator-validators

Execute the DelegatorValidators RPC method

```
axoned query distribution delegator-validators [flags]
```

### Options

```
      --delegator-address account address or key name   
      --grpc-addr string                                the gRPC endpoint to use for this chain
      --grpc-insecure                                   allow gRPC over insecure channels, if not the server must use TLS
      --height int                                      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                                            help for delegator-validators
      --no-indent                                       Do not indent JSON output
      --node string                                     <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string                                   Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query distribution](axoned_query_distribution.md)	 - Querying commands for the distribution module
