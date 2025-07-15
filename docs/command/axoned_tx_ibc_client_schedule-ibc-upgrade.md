## axoned tx ibc client schedule-ibc-upgrade

Submit an IBC software upgrade proposal

### Synopsis

Please specify a unique name and height for the upgrade to take effect.
		The client state specified is the upgraded client state representing the upgraded chain
		
		Example Upgraded Client State JSON: 
		\{
			"@type":"/ibc.lightclients.tendermint.v1.ClientState",
			"chain_id":"testchain1",
			"unbonding_period":"1814400s",
			"latest_height":\{
			   "revision_number":"0",
			   "revision_height":"2"
			\},
			"proof_specs":[
			   \{
				  "leaf_spec":\{
					 "hash":"SHA256",
					 "prehash_key":"NO_HASH",
					 "prehash_value":"SHA256",
					 "length":"VAR_PROTO",
					 "prefix":"AA=="
				  \},
				  "inner_spec":\{
					 "child_order":[
						0,
						1
					 ],
					 "child_size":33,
					 "min_prefix_length":4,
					 "max_prefix_length":12,
					 "empty_child":null,
					 "hash":"SHA256"
				  \},
				  "max_depth":0,
				  "min_depth":0
			   \},
			   \{
				  "leaf_spec":\{
					 "hash":"SHA256",
					 "prehash_key":"NO_HASH",
					 "prehash_value":"SHA256",
					 "length":"VAR_PROTO",
					 "prefix":"AA=="
				  \},
				  "inner_spec":\{
					 "child_order":[
						0,
						1
					 ],
					 "child_size":32,
					 "min_prefix_length":1,
					 "max_prefix_length":1,
					 "empty_child":null,
					 "hash":"SHA256"
				  \},
				  "max_depth":0,
				  "min_depth":0
			   \}
			],
			"upgrade_path":[
			   "upgrade",
			   "upgradedIBCState"
			]
		 \}
		
```
axoned tx ibc client schedule-ibc-upgrade [name] [height] [path/to/upgraded_client_state.json] [flags]
```

### Options

```
  -a, --account-number uint      The account number of the signing account (offline mode only)
      --authority string         The address of the client module authority (defaults to gov)
      --aux                      Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string    Transaction broadcasting mode (sync|async) (default "sync")
      --chain-id string          The network chain ID
      --deposit string           The deposit to include with the governance proposal
      --dry-run                  ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string       Fee granter grants fees for the transaction
      --fee-payer string         Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string              Fees to pay along with transaction; eg: 10uatom
      --from string              Name or address of private key with which to sign
      --gas string               gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float     adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string        Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only            Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                     help for schedule-ibc-upgrade
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) (default "os")
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                   Use a connected Ledger device
      --metadata string          The metadata to include with the governance proposal
      --node string              <host>:<port> to CometBFT rpc interface for this chain (default "tcp://localhost:26657")
      --note string              Note to add a description to the transaction (previously --memo)
      --offline                  Offline mode (does not allow any online functionality)
  -o, --output string            Output format (text|json) (default "json")
  -s, --sequence uint            The sequence number of the signing account (offline mode only)
      --sign-mode string         Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --summary string           The summary to include with the governance proposal
      --timeout-height uint      Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string               Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --title string             The title to put on the governance proposal
  -y, --yes                      Skip tx broadcasting prompt confirmation
```

### SEE ALSO

* [axoned tx ibc client](axoned_tx_ibc_client.md)	 - IBC client transaction subcommands
