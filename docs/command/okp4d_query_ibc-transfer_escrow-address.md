## okp4d query ibc-transfer escrow-address

Get the escrow address for a channel

### Synopsis

Get the escrow address for a channel

```
okp4d query ibc-transfer escrow-address [flags]
```

### Examples

```
okp4d query ibc-transfer escrow-address [port] [channel-id]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for escrow-address
      --node string     &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query ibc-transfer](okp4d_query_ibc-transfer.md)	 - IBC fungible token transfer query subcommands

