## axoned query block

Query for a committed block by height, hash, or event(s)

### Synopsis

Query for a specific committed block using the CometBFT RPC `block` and `block_by_hash` method

```
axoned query block --type=[height|hash] [height|hash] [flags]
```

### Examples

```
$ axoned query block --type=height <height>
$ axoned query block --type=hash <hash>
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for block
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
      --type string        The type to be used when querying tx, can be one of "height", "hash" (default "hash")
```

### SEE ALSO

* [axoned query](axoned_query.md)	 - Querying subcommands
