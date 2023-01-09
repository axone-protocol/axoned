## okp4d query logic ask

executes a logic query and returns the solutions found.

### Synopsis

Executes the [query] for the given [program] file and return the solution(s) found.

Since the query is without any side-effect, the query is not executed in the context of a transaction and no fee
is charged for this, but the execution is constrained by the current limits configured in the module.

Example:
$ okp4d logic query ask "immortal(X)." program.txt

```
okp4d query logic ask [query] [program] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for ask
      --node string        &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query logic](okp4d_query_logic.md)	 - Querying commands for the logic module
