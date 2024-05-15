## axoned query slashing signing-info

Query a validator's signing information

### Synopsis

Query a validator's signing information, with a pubkey ('axoned comet show-validator') or a validator consensus address

```
axoned query slashing signing-info [validator-conspub/address] [flags]
```

### Examples

```
axoned query slashing signing-info '{"@type":"/cosmos.crypto.ed25519.PubKey","key":"OauFcTKbN5Lx3fJL689cikXBqe+hcp6Y+x0rYUdR9Jk="}'
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for signing-info
      --no-indent          Do not indent JSON output
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query slashing](axoned_query_slashing.md)	 - Querying commands for the slashing module
