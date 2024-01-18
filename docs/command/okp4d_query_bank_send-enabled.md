## okp4d query bank send-enabled

Query for send enabled entries

### Synopsis

Query for send enabled entries that have been specifically set.

To look up one or more specific denoms, supply them as arguments to this command.
To look up all denoms, do not provide any arguments.

```
okp4d query bank send-enabled [denom1 ...] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for send-enabled
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

* [okp4d query bank](okp4d_query_bank.md)	 - Querying commands for the bank module
