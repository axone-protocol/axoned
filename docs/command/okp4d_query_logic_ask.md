## okp4d query logic ask

Executes a logic query and returns the solution(s) found.

### Synopsis

Executes the [query] and return the solution(s) found.

Optionally, a program can be transmitted, which will be compiled before the query is processed.

Since the query is without any side-effect, the query is not executed in the context of a transaction and no fee
is charged for this, but the execution is constrained by the current limits configured in the module (that you can
query).

```
okp4d query logic ask [query] [flags]
```

### Examples

```
$ okp4d query logic ask "chain_id(X)." # returns the chain-id
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for ask
      --no-indent          Do not indent JSON output
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
  -p, --program string     The program to compile before the query.
```

### SEE ALSO

* [okp4d query logic](okp4d_query_logic.md)	 - Querying commands for the logic module
