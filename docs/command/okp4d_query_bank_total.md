## okp4d query bank total

Query the total supply of coins of the chain

### Synopsis

Query total supply of coins that are held by accounts in the chain.

Example:
  $ okp4d query bank total

To query for the total supply of a specific coin denomination use:
  $ okp4d query bank total --denom=[denom]

```
okp4d query bank total [flags]
```

### Options

```
      --count-total       count total number of records in all supply totals to query for
      --denom string      The specific balance denomination to query for
      --height int        Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help              help for total
      --limit uint        pagination limit of all supply totals to query for (default 100)
      --node string       &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --offset uint       pagination offset of all supply totals to query for
  -o, --output string     Output format (text|json) (default "text")
      --page uint         pagination page of all supply totals to query for. This sets offset to a multiple of limit (default 1)
      --page-key string   pagination page-key of all supply totals to query for
      --reverse           results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query bank](okp4d_query_bank.md)	 - Querying commands for the bank module

