## axoned query consensus comet block-by-height

Query for a committed block by height

### Synopsis

Query for a specific committed block using the CometBFT RPC `block_by_height` method

```
axoned query consensus comet block-by-height [height] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for block-by-height
      --no-indent          Do not indent JSON output
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query consensus comet](axoned_query_consensus_comet.md)	 - Querying commands for the cosmos.base.tendermint.v1beta1.Service service
