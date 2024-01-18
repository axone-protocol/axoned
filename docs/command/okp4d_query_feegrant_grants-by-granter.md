## okp4d query feegrant grants-by-granter

Query all grants by a granter

```
okp4d query feegrant grants-by-granter [granter] [flags]
```

### Examples

```
$ okp4d query feegrant grants-by-granter [granter]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for grants-by-granter
      --no-indent          Do not indent JSON output
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
      --page-count-total   
      --page-key binary    
      --page-limit uint    
      --page-offset uint   
      --page-reverse       
```

### SEE ALSO

* [okp4d query feegrant](okp4d_query_feegrant.md)	 - Querying commands for the feegrant module
