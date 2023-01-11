## okp4d query distribution rewards

Query all distribution delegator rewards or rewards from a particular validator

### Synopsis

Query all rewards earned by a delegator, optionally restrict to rewards from a single validator.

Example:
$ okp4d query distribution rewards okp41gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
$ okp4d query distribution rewards okp41gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p okp4valoper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj

```
okp4d query distribution rewards [delegator-addr] [validator-addr] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for rewards
      --node string        &lt;host&gt;:&lt;port&gt; to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d query distribution](okp4d_query_distribution.md)	 - Querying commands for the distribution module
