## okp4d query bank denom-metadata

Query the client metadata for coin denominations

### Synopsis

Query the client metadata for all the registered coin denominations

Example:
  To query for the client metadata of all coin denominations use:
  $ okp4d query bank denom-metadata

To query for the client metadata of a specific coin denomination use:
  $ okp4d query bank denom-metadata --denom=[denom]

```
okp4d query bank denom-metadata [flags]
```

### Options

```
      --denom string       The specific denomination to query client metadata for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for denom-metadata
      --node string        &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query bank](okp4d_query_bank.md)	 - Querying commands for the bank module
