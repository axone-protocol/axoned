## okp4d query ibc client consensus-state

Query the consensus state of a client at a given height

### Synopsis

Query the consensus state for a particular light client at a given height.
If the '--latest' flag is included, the query returns the latest consensus state, overriding the height argument.

```
okp4d query ibc client consensus-state [client-id] [height] [flags]
```

### Examples

```
okp4d query ibc client  consensus-state [client-id] [height]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for consensus-state
      --latest-height   return latest stored consensus state
      --node string     &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
      --prove           show proofs for the query results (default true)
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query ibc client](okp4d_query_ibc_client.md)	 - IBC client query subcommands
