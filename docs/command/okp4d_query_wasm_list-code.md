## okp4d query wasm list-code

List all wasm bytecode on the chain

### Synopsis

List all wasm bytecode on the chain

```
okp4d query wasm list-code [flags]
```

### Options

```
      --count-total        count total number of records in list codes to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-code
      --limit uint         pagination limit of list codes to query for (default 100)
      --node string        &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --offset uint        pagination offset of list codes to query for
  -o, --output string      Output format (text|json) (default "text")
      --page uint          pagination page of list codes to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of list codes to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query wasm](okp4d_query_wasm.md)	 - Querying commands for the wasm module
