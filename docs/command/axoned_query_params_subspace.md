## axoned query params subspace

Query for raw parameters by subspace and key

```
axoned query params subspace [subspace] [key] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for subspace
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) (default "os")
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string            Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query params](axoned_query_params.md)	 - Querying commands for the params module
