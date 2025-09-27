## axoned keys import-hex

Import private keys into the local keybase

### Synopsis

Import hex encoded private key into the local keybase.
Supported key-types can be obtained with:
axoned list-key-types

```
axoned keys import-hex <name> [hex] [flags]
```

### Options

```
  -h, --help              help for import-hex
      --key-type string   private key signing algorithm kind (default "secp256k1")
```

### Options inherited from parent commands

```
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) (default "test")
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --output string            Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned keys](axoned_keys.md)	 - Manage your application's keys
