## okp4d query tx

Query for a transaction by hash, "&lt;addr&gt;/&lt;seq&gt;" combination or comma-separated signatures in a committed block

### Synopsis

Example:
$ okp4d query tx &lt;hash&gt;
$ okp4d query tx --type=acc_seq &lt;addr&gt;/&lt;sequence&gt;
$ okp4d query tx --type=signature &lt;sig1_base64&gt;,&lt;sig2_base64...&gt;

```
okp4d query tx --type=[hash|acc_seq|signature] [hash|acc_seq|signature] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for tx
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
      --type string        The type to be used when querying tx, can be one of "hash", "acc_seq", "signature" (default "hash")
```

### SEE ALSO

* [okp4d query](okp4d_query.md)	 - Querying subcommands
