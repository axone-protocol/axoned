## axoned query logic program

shows stored program metadata

```
axoned query logic program [program-id] [flags]
```

### Examples

```
$ axoned query logic program 0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for program
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query logic](axoned_query_logic.md)	 - Querying commands for the logic module
