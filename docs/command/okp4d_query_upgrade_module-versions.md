## okp4d query upgrade module-versions

Query the list of module versions

### Synopsis

Gets a list of module names and their respective consensus versions. Following the command with a specific module name will return only that module's information.

```
okp4d query upgrade module-versions [optional module_name] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for module-versions
      --no-indent          Do not indent JSON output
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [okp4d query upgrade](okp4d_query_upgrade.md)	 - Querying commands for the upgrade module
