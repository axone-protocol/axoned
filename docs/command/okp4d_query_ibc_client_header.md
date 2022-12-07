## okp4d query ibc client header

Query the latest header of the running chain

### Synopsis

Query the latest Tendermint header of the running chain

```
okp4d query ibc client header [flags]
```

### Examples

```
okp4d query ibc client  header
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for header
      --node string     &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query ibc client](okp4d_query_ibc_client.md)	 - IBC client query subcommands

