## axoned query staking unbonding-delegation

Query an unbonding-delegation record based on delegator and validator address

### Synopsis

Query unbonding delegations for an individual delegator on an individual validator.

```
axoned query staking unbonding-delegation [delegator-addr] [validator-addr] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for unbonding-delegation
      --no-indent          Do not indent JSON output
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query staking](axoned_query_staking.md)	 - Querying commands for the staking module
