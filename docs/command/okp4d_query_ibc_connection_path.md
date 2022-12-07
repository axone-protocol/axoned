## okp4d query ibc connection path

Query stored client connection paths

### Synopsis

Query stored client connection paths

```
okp4d query ibc connection path [client-id] [flags]
```

### Examples

```
okp4d query  ibc connection path [client-id]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for path
      --node string     &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
      --prove           show proofs for the query results (default true)
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query ibc connection](okp4d_query_ibc_connection.md)	 - IBC connection query subcommands

