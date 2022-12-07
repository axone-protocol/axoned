## okp4d query ibc-fee payee

Query the relayer payee address on a given channel

### Synopsis

Query the relayer payee address on a given channel

```
okp4d query ibc-fee payee [channel-id] [relayer] [flags]
```

### Examples

```
okp4d query ibc-fee payee channel-5 cosmos1layxcsmyye0dc0har9sdfzwckaz8sjwlfsj8zs
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for payee
      --node string     &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query ibc-fee](okp4d_query_ibc-fee.md)	 - IBC relayer incentivization query subcommands

