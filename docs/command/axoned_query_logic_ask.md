## axoned query logic ask

executes a logic query and returns the solutions found.

### Synopsis

Executes the [query] and return the solution(s) found.
 Optionally, a program can be transmitted, which will be interpreted before the query is processed.
 Since the query is without any side-effect, the query is not executed in the context of a transaction and no fee
 is charged for this, but the execution is constrained by the current limits configured in the module (that you can
 query).

```
axoned query logic ask [query] [flags]
```

### Examples

```
$ axoned query logic ask "chain_id(X)." # returns the chain-id
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for ask
      --limit uint         limit the maximum number of solutions to return.
                           This parameter is constrained by the 'max_result_count' setting in the module configuration, which specifies the maximum number of results that can be requested per query. (default 1)
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
      --program string     reads the program from the given string.
```

### SEE ALSO

* [axoned query logic](axoned_query_logic.md)	 - Querying commands for the logic module
