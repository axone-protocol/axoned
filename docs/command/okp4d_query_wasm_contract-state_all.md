## okp4d query wasm contract-state all

Prints out all internal state of a contract given its address

### Synopsis

Prints out all internal state of a contract given its address

```
okp4d query wasm contract-state all [bech32_address] [flags]
```

### Options

```
      --count-total        count total number of records in contract state to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for all
      --limit uint         pagination limit of contract state to query for (default 100)
      --node string        &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --offset uint        pagination offset of contract state to query for
  -o, --output string      Output format (text|json) (default "text")
      --page uint          pagination page of contract state to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of contract state to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query wasm contract-state](okp4d_query_wasm_contract-state.md)	 - Querying commands for the wasm module
