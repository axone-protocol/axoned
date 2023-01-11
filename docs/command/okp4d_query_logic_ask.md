## okp4d query logic ask

executes a logic query and returns the solutions found.

### Synopsis

Executes the [query] and return the solution(s) found.

Optionally, a program can be transmitted, which will be interpreted before the query is processed.

Since the query is without any side-effect, the query is not executed in the context of a transaction and no fee
is charged for this, but the execution is constrained by the current limits configured in the module (that you can
query).

```
okp4d query logic ask [query] [flags]
```

### Examples

```
okp4d query logic ask "chain_id(X)." # returns the chain-id
```

### Options

```
      --grpc-addr string      the gRPC endpoint to use for this chain
      --grpc-insecure         allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int            Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                  help for ask
      --node string           &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string         Output format (text|json) (default "text")
      --program string        reads the program from the given filename or from stdin if "-" is passed as the filename.
      --program-file string   reads the program from the given string.
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query logic](okp4d_query_logic.md)	 - Querying commands for the logic module
