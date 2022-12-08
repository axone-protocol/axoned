## okp4d query feegrant grant

Query details of a single grant

### Synopsis

Query details for a grant.
You can find the fee-grant of a granter and grantee.

Example:
$ okp4d query feegrant grant [granter] [grantee]

```
okp4d query feegrant grant [granter] [grantee] [flags]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for grant
      --node string     &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query feegrant](okp4d_query_feegrant.md)	 - Querying commands for the feegrant module
