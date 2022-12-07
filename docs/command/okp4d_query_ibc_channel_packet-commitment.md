## okp4d query ibc channel packet-commitment

Query a packet commitment

### Synopsis

Query a packet commitment

```
okp4d query ibc channel packet-commitment [port-id] [channel-id] [sequence] [flags]
```

### Examples

```
okp4d query ibc channel packet-commitment [port-id] [channel-id] [sequence]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for packet-commitment
      --node string     &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
      --prove           show proofs for the query results (default true)
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query ibc channel](okp4d_query_ibc_channel.md)	 - IBC channel query subcommands

