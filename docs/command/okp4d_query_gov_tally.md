## okp4d query gov tally

Get the tally of a proposal vote

### Synopsis

Query tally of votes on a proposal. You can find
the proposal-id by running "okp4d query gov proposals".

Example:
$ okp4d query gov tally 1

```
okp4d query gov tally [proposal-id] [flags]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for tally
      --node string     &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query gov](okp4d_query_gov.md)	 - Querying commands for the governance module
