## okp4d query upgrade plan

get upgrade plan (if one exists)

### Synopsis

Gets the currently scheduled upgrade plan, if one exists

```
okp4d query upgrade plan [flags]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for plan
      --node string     &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query upgrade](okp4d_query_upgrade.md)	 - Querying commands for the upgrade module
