## okp4d query group group-policies-by-group

Query for group policies by group id with pagination flags

```
okp4d query group group-policies-by-group [group-id] [flags]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for group-policies-by-group
      --node string     <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query group](okp4d_query_group.md)	 - Querying commands for the group module
