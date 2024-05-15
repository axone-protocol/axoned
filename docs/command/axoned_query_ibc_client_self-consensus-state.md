## axoned query ibc client self-consensus-state

Query the self consensus state for this chain

### Synopsis

Query the self consensus state for this chain. This result may be used for verifying IBC clients representing this chain which are hosted on counterparty chains.

```
axoned query ibc client self-consensus-state [flags]
```

### Examples

```
axoned query ibc client self-consensus-state
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for self-consensus-state
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query ibc client](axoned_query_ibc_client.md)	 - IBC client query subcommands
