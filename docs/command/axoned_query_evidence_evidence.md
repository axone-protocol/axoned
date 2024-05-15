## axoned query evidence evidence

Query for evidence by hash

```
axoned query evidence evidence [hash] [flags]
```

### Examples

```
axoned query evidence DF0C23E8634E480F84B9D5674A7CDC9816466DEC28A3358F73260F68D28D7660
```

### Options

```
      --evidence-hash binary   
      --grpc-addr string       the gRPC endpoint to use for this chain
      --grpc-insecure          allow gRPC over insecure channels, if not the server must use TLS
      --height int             Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                   help for evidence
      --no-indent              Do not indent JSON output
      --node string            <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string          Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned query evidence](axoned_query_evidence.md)	 - Querying commands for the evidence module
