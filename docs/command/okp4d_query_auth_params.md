## okp4d query auth params

Query the current auth parameters

### Synopsis

Query the current auth parameters:

$ okp4d query auth params

```
okp4d query auth params [flags]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for params
      --node string     &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query auth](okp4d_query_auth.md)	 - Querying commands for the auth module

