## axoned query distribution validator-historical-rewards

Query validator historical rewards for a specific period

```
axoned query distribution validator-historical-rewards [validator] [period] [flags]
```

### Examples

```
$ axoned query distribution validator-historical-rewards [validator-address] 5
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for validator-historical-rewards
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) (default "os")
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string            Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query distribution](axoned_query_distribution.md)	 - Querying commands for the distribution module
