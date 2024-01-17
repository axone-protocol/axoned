## okp4d query ibc-fee counterparty-payee

Query the relayer counterparty payee on a given channel

### Synopsis

Query the relayer counterparty payee on a given channel

```
okp4d query ibc-fee counterparty-payee [channel-id] [relayer] [flags]
```

### Examples

```
okp4d query ibc-fee counterparty-payee channel-5 cosmos1layxcsmyye0dc0har9sdfzwckaz8sjwlfsj8zs
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for counterparty-payee
      --node string        <host>:<port> to CometBFT RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string      Output format (text|json) (default "text")
```

### SEE ALSO

* [okp4d query ibc-fee](okp4d_query_ibc-fee.md)	 - IBC relayer incentivization query subcommands
