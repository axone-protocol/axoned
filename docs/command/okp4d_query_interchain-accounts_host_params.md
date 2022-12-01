## okp4d query interchain-accounts host params

Query the current interchain-accounts host submodule parameters

### Synopsis

Query the current interchain-accounts host submodule parameters

```
okp4d query interchain-accounts host params [flags]
```

### Examples

```
<appd> query interchain-accounts host params
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for params
      --node string     <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query interchain-accounts host](okp4d_query_interchain-accounts_host.md)	 - interchain-accounts host subcommands

