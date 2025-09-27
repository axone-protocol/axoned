## axoned keys show

Retrieve key information by name or address

### Synopsis

Display keys details. If multiple names or addresses are provided,
then an ephemeral multisig key will be created under the name "multi"
consisting of all the keys provided by name and multisig threshold.

```
axoned keys show [name_or_address [name_or_address...]] [flags]
```

### Options

```
  -a, --address                  Output the address only (cannot be used with --output)
      --bech string              The Bech32 prefix encoding for a key (acc|val|cons) (default "acc")
  -d, --device                   Output the address in a ledger device (cannot be used with --pubkey)
  -h, --help                     help for show
      --multisig-threshold int   K out of N required signatures (default 1)
  -p, --pubkey                   Output the public key only (cannot be used with --output)
      --qrcode                   Display key address QR code (will be ignored if -a or --address is false)
```

### Options inherited from parent commands

```
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) (default "test")
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --output string            Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned keys](axoned_keys.md)	 - Manage your application's keys
