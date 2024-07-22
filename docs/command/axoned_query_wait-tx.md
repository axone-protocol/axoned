## axoned query wait-tx

Wait for a transaction to be included in a block

### Synopsis

Subscribes to a CometBFT WebSocket connection and waits for a transaction event with the given hash.

```
axoned query wait-tx [hash] [flags]
```

### Examples

```
By providing the transaction hash:
$ axoned q wait-tx [hash]

Or, by piping a "tx" command:
$ axoned tx [flags] | axoned q wait-tx

```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for wait-tx
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
      --timeout duration   The maximum time to wait for the transaction to be included in a block (default 15s)
```

### SEE ALSO

* [axoned query](axoned_query.md)	 - Querying subcommands
