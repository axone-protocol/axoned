## okp4d query distribution delegation-total-rewards

Execute the DelegationTotalRewards RPC method

```
okp4d query distribution delegation-total-rewards [flags]
```

### Options

```
      --delegator-address account address or key name   
      --grpc-addr string                                the gRPC endpoint to use for this chain
      --grpc-insecure                                   allow gRPC over insecure channels, if not the server must use TLS
      --height int                                      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                                            help for delegation-total-rewards
      --no-indent                                       Do not indent JSON output
      --node string                                     <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string                                   Output format (text|json) (default "text")
```

### SEE ALSO

* [okp4d query distribution](okp4d_query_distribution.md)	 - Querying commands for the distribution module
