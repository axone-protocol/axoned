## okp4d query gov proposer

Query the proposer of a governance proposal

### Synopsis

Query which address proposed a proposal with a given ID

```
okp4d query gov proposer [proposal-id] [flags]
```

### Examples

```
okp4d query gov proposer 1
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for proposer
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [okp4d query gov](okp4d_query_gov.md)	 - Querying commands for the gov module
