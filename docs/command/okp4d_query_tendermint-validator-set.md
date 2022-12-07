## okp4d query tendermint-validator-set

Get the full tendermint validator set at given height

```
okp4d query tendermint-validator-set [height] [flags]
```

### Options

```
  -h, --help            help for tendermint-validator-set
      --limit int       Query number of results returned per page (default 100)
      --node string     &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
      --page int        Query a specific page of paginated results (default 1)
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query](okp4d_query.md)	 - Querying subcommands
