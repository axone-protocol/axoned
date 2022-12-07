## okp4d query interchain-accounts host packet-events

Query the interchain-accounts host submodule packet events

### Synopsis

Query the interchain-accounts host submodule packet events for a particular channel and sequence

```
okp4d query interchain-accounts host packet-events [channel-id] [sequence] [flags]
```

### Examples

```
okp4d query interchain-accounts host packet-events channel-0 100
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for packet-events
      --node string     &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query interchain-accounts host](okp4d_query_interchain-accounts_host.md)	 - interchain-accounts host subcommands
