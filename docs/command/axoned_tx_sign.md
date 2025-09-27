## axoned tx sign

Sign a transaction generated offline

### Synopsis

Sign a transaction created with the --generate-only flag.
It will read a transaction from [file], sign it, and print its JSON encoding.

If the --signature-only flag is set, it will output the signature parts only.

The --offline flag makes sure that the client will not reach out to full node.
As a result, the account and sequence number queries will not be performed and
it is required to set such parameters manually. Note, invalid values will cause
the transaction to fail.

The --multisig=&lt;multisig_key&gt; flag generates a signature on behalf of a multisig account
key. It implies --signature-only. Full multisig signed transactions may eventually
be generated via the 'multisign' command.

```
axoned tx sign [file] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) (default "sync")
      --chain-id string             The network chain ID (default "axoned")
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for sign
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) (default "test")
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --multisig string             Address or key name of the multisig account on behalf of which the transaction shall be signed
      --node string                 <host>:<port> to CometBFT rpc interface for this chain (default "tcp://localhost:26657")
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) (default "json")
      --output-document string      The document will be written to the given file instead of STDOUT
      --overwrite                   Overwrite existing signatures with a new one. If disabled, new signature will be appended
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --signature-only              Print only the signatures
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
```

### SEE ALSO

* [axoned tx](axoned_tx.md)	 - Transactions subcommands
