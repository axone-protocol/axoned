## okp4d query gov params

Query the parameters of the governance process

### Synopsis

Query the all the parameters for the governance process.

Example:
$ <appd> query gov params

```
okp4d query gov params [flags]
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

* [okp4d query gov](okp4d_query_gov.md)	 - Querying commands for the governance module
