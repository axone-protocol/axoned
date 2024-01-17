## okp4d query auth module-account

Query module account info by module name

```
okp4d query auth module-account [module-name] [flags]
```

### Examples

```
okp4d q auth module-account gov
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for module-account
      --no-indent          Do not indent JSON output
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [okp4d query auth](okp4d_query_auth.md)	 - Querying commands for the auth module
