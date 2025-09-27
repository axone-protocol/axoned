## axoned tx multi-sign

Generate multisig signatures for transactions generated offline

### Synopsis

Sign transactions created with the --generate-only flag that require multisig signatures.

Read one or more signatures from one or more [signature] file, generate a multisig signature compliant to the
multisig key [name], and attach the key name to the transaction read from [file].

Example:
$ axoned tx multisign transaction.json k1k2k3 k1sig.json k2sig.json k3sig.json

If --signature-only flag is on, output a JSON representation
of only the generated signature.

If the --offline flag is on, the client will not reach out to an external node.
Account number or sequence number lookups are not performed so you must
set these parameters manually.

If the --skip-signature-verification flag is on, the command will not verify the
signatures in the provided signature files. This is useful when the multisig
account is a signer in a nested multisig scenario.

The current multisig implementation defaults to amino-json sign mode.
The SIGN_MODE_DIRECT sign mode is not supported.'

```
axoned tx multi-sign [file] [name] [[signature]...] [flags]
```

### Options

```
  -a, --account-number uint           The account number of the signing account (offline mode only)
      --aux                           Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string         Transaction broadcasting mode (sync|async) (default "sync")
      --chain-id string               The network chain ID (default "axoned")
      --dry-run                       ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string            Fee granter grants fees for the transaction
      --fee-payer string              Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                   Fees to pay along with transaction; eg: 10uatom
      --from string                   Name or address of private key with which to sign
      --gas string                    gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float          adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string             Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only                 Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                          help for multi-sign
      --keyring-backend string        Select keyring's backend (os|file|kwallet|pass|test|memory) (default "test")
      --keyring-dir string            The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                        Use a connected Ledger device
      --node string                   <host>:<port> to CometBFT rpc interface for this chain (default "tcp://localhost:26657")
      --note string                   Note to add a description to the transaction (previously --memo)
      --offline                       Offline mode (does not allow any online functionality)
      --output-document string        The document is written to the given file instead of STDOUT
  -s, --sequence uint                 The sequence number of the signing account (offline mode only)
      --sign-mode string              Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --signature-only                Print only the generated signature, then exit
      --skip-signature-verification   Skip signature verification
      --timeout-duration duration     TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint           DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                    Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                     Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                           Skip tx broadcasting prompt confirmation
```

### SEE ALSO

* [axoned tx](axoned_tx.md)	 - Transactions subcommands
