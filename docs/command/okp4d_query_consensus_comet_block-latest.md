## okp4d query consensus comet block-latest

Query for the latest committed block

```
okp4d query consensus comet block-latest [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for block-latest
      --no-indent          Do not indent JSON output
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [okp4d query consensus comet](okp4d_query_consensus_comet.md)	 - Querying commands for the cosmos.base.tendermint.v1beta1.Service service
