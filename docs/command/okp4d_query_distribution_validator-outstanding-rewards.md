## okp4d query distribution validator-outstanding-rewards

Query distribution outstanding (un-withdrawn) rewards for a validator and all their delegations

### Synopsis

Query distribution outstanding (un-withdrawn) rewards for a validator and all their delegations.

Example:
$ okp4d query distribution validator-outstanding-rewards okp4valoper1lwjmdnks33xwnmfayc64ycprww49n33mtm92ne

```
okp4d query distribution validator-outstanding-rewards [validator] [flags]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for validator-outstanding-rewards
      --node string     &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query distribution](okp4d_query_distribution.md)	 - Querying commands for the distribution module

