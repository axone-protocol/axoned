## okp4d query ibc client self-consensus-state

Query the self consensus state for this chain

### Synopsis

Query the self consensus state for this chain. This result may be used for verifying IBC clients representing this chain which are hosted on counterparty chains.

```
okp4d query ibc client self-consensus-state [flags]
```

### Examples

```
okp4d query ibc client self-consensus-state
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for self-consensus-state
      --node string     &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query ibc client](okp4d_query_ibc_client.md)	 - IBC client query subcommands

