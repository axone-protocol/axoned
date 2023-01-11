## okp4d query auth address-by-acc-num

Query for an address by account number

```
okp4d query auth address-by-acc-num [acc-num] [flags]
```

### Examples

```
okp4d q auth address-by-acc-num 1
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for address-by-acc-num
      --node string        &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query auth](okp4d_query_auth.md)	 - Querying commands for the auth module
