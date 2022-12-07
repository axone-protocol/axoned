## okp4d query wasm code

Downloads wasm bytecode for given code id

### Synopsis

Downloads wasm bytecode for given code id

```
okp4d query wasm code [code_id] [output filename] [flags]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for code
      --node string     &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query wasm](okp4d_query_wasm.md)	 - Querying commands for the wasm module

