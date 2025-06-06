## axoned tx wasm submit-proposal instantiate-contract-2

Submit an instantiate wasm contract proposal with predictable address

```
axoned tx wasm submit-proposal instantiate-contract-2 [code_id_int64] [json_encoded_init_args] [salt] --authority [address] --label [text] --title [text] --summary [text] --admin [address,optional] --amount [coins,optional] --fix-msg [bool,optional] [flags]
```

### Options

```
  -a, --account-number uint      The account number of the signing account (offline mode only)
      --admin string             Address of an admin
      --amount string            Coins to send to the contract during instantiation
      --ascii                    ascii encoded salt
      --authority string         The address of the governance account. Default is the sdk gov module account (default "axone10d07y265gmmuvt4z0w9aw880jnsr700jeyljzw")
      --aux                      Generate aux signer data instead of sending a tx
      --b64                      base64 encoded salt
  -b, --broadcast-mode string    Transaction broadcasting mode (sync|async) (default "sync")
      --chain-id string          The network chain ID
      --deposit string           Deposit of proposal
      --dry-run                  ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --expedite                 Expedite proposals have shorter voting period but require higher voting threshold
      --fee-granter string       Fee granter grants fees for the transaction
      --fee-payer string         Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string              Fees to pay along with transaction; eg: 10uatom
      --fix-msg                  An optional flag to include the json_encoded_init_args for the predictable address generation mode
      --from string              Name or address of private key with which to sign
      --gas string               gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float     adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string        Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only            Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                     help for instantiate-contract-2
      --hex                      hex encoded salt
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) (default "os")
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --label string             A human-readable name for this contract in lists
      --ledger                   Use a connected Ledger device
      --no-admin                 You must set this explicitly if you don't want an admin
      --node string              <host>:<port> to CometBFT rpc interface for this chain (default "tcp://localhost:26657")
      --note string              Note to add a description to the transaction (previously --memo)
      --offline                  Offline mode (does not allow any online functionality)
  -o, --output string            Output format (text|json) (default "json")
  -s, --sequence uint            The sequence number of the signing account (offline mode only)
      --sign-mode string         Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --summary string           Summary of proposal
      --timeout-height uint      Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string               Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --title string             Title of proposal
  -y, --yes                      Skip tx broadcasting prompt confirmation
```

### SEE ALSO

* [axoned tx wasm submit-proposal](axoned_tx_wasm_submit-proposal.md)	 - Submit a wasm proposal.
