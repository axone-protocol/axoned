## okp4d query wasm contract-state raw

Prints out internal state for key of a contract given its address

### Synopsis

Prints out internal state for of a contract given its address

```
okp4d query wasm contract-state raw [bech32_address] [key] [flags]
```

### Options

```
      --ascii           ascii encoded key argument
      --b64             base64 encoded key argument
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for raw
      --hex             hex encoded  key argument
      --node string     <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query wasm contract-state](okp4d_query_wasm_contract-state.md)	 - Querying commands for the wasm module

