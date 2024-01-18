## okp4d query wasm contract

Prints out metadata of a contract given its address

### Synopsis

Prints out metadata of a contract given its address

```
okp4d query wasm contract [bech32_address] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for contract
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [okp4d query wasm](okp4d_query_wasm.md)	 - Querying commands for the wasm module
