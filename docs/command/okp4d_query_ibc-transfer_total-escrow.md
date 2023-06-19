## okp4d query ibc-transfer total-escrow

Query the total amount of tokens in escrow for a denom

### Synopsis

Query the total amount of tokens in escrow for a denom

```
okp4d query ibc-transfer total-escrow [denom] [flags]
```

### Examples

```
okp4d query ibc-transfer total-escrow uosmo
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for total-escrow
      --node string        &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query ibc-transfer](okp4d_query_ibc-transfer.md)	 - IBC fungible token transfer query subcommands
