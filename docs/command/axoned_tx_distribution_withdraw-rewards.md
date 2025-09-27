## axoned tx distribution withdraw-rewards

Withdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator

### Synopsis

Withdraw rewards from a given delegation address,
and optionally withdraw validator commission if the delegation address given is a validator operator.

Example:
$ axoned tx distribution withdraw-rewards axonevaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj --from mykey
$ axoned tx distribution withdraw-rewards axonevaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj --from mykey --commission

```
axoned tx distribution withdraw-rewards [validator-addr] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) (default "sync")
      --chain-id string             The network chain ID
      --commission                  Withdraw the validator's commission in addition to the rewards
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for withdraw-rewards
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) (default "os")
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 <host>:<port> to CometBFT rpc interface for this chain (default "tcp://localhost:26657")
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) (default "json")
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
```

### SEE ALSO

* [axoned tx distribution](axoned_tx_distribution.md)	 - Distribution transactions subcommands
