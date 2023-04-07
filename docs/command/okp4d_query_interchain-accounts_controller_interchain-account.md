## okp4d query interchain-accounts controller interchain-account

Query the interchain account address for a given owner on a particular connection

### Synopsis

Query the controller submodule for the interchain account address for a given owner on a particular connection

```
okp4d query interchain-accounts controller interchain-account [owner] [connection-id] [flags]
```

### Examples

```
okp4d query interchain-accounts controller interchain-account cosmos1layxcsmyye0dc0har9sdfzwckaz8sjwlfsj8zs connection-0
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for interchain-account
      --node string        &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query interchain-accounts controller](okp4d_query_interchain-accounts_controller.md)	 - IBC interchain accounts controller query subcommands
