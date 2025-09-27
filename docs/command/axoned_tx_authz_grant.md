## axoned tx authz grant

Grant authorization to an address

### Synopsis

create a new grant authorization to an address to execute a transaction on your behalf:

Examples:
 $ axoned tx authz grant cosmos1skjw.. send --spend-limit=1000stake --from=cosmos1skl..
 $ axoned tx authz grant cosmos1skjw.. generic --msg-type=/cosmos.gov.v1.MsgVote --from=cosmos1sk..

```
axoned tx authz grant <grantee> <authorization_type="send"|"generic"|"delegate"|"unbond"|"redelegate"> --from <granter> [flags]
```

### Options

```
  -a, --account-number uint          The account number of the signing account (offline mode only)
      --allow-list strings           Allowed addresses grantee is allowed to send funds separated by ,
      --allowed-validators strings   Allowed validators addresses separated by ,
      --aux                          Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string        Transaction broadcasting mode (sync|async) (default "sync")
      --chain-id string              The network chain ID
      --deny-validators strings      Deny validators addresses separated by ,
      --dry-run                      ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --expiration int               Expire time as Unix timestamp. Set zero (0) for no expiry. Default is 0.
      --fee-granter string           Fee granter grants fees for the transaction
      --fee-payer string             Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                  Fees to pay along with transaction; eg: 10uatom
      --from string                  Name or address of private key with which to sign
      --gas string                   gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float         adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string            Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only                Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                         help for grant
      --keyring-backend string       Select keyring's backend (os|file|kwallet|pass|test|memory) (default "os")
      --keyring-dir string           The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                       Use a connected Ledger device
      --msg-type string              The Msg method name for which we are creating a GenericAuthorization
      --node string                  <host>:<port> to CometBFT rpc interface for this chain (default "tcp://localhost:26657")
      --note string                  Note to add a description to the transaction (previously --memo)
      --offline                      Offline mode (does not allow any online functionality)
  -o, --output string                Output format (text|json) (default "json")
  -s, --sequence uint                The sequence number of the signing account (offline mode only)
      --sign-mode string             Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --spend-limit string           SpendLimit for Send Authorization, an array of Coins allowed spend
      --timeout-duration duration    TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint          DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                   Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                    Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                          Skip tx broadcasting prompt confirmation
```

### SEE ALSO

* [axoned tx authz](axoned_tx_authz.md)	 - Authorization transactions subcommands
