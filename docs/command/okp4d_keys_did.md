## okp4d keys did

Give the did:key from a ed25519 or secp256k1 pubkey (hex, base64)

### Synopsis

Give the did:key from a ed25519 or secp256k1 pubkey given as hex or base64 encoded string.

Example:

    $ okp4d keys did "AtD+mbIUqu615Grk1loWI6ldnQzs1X1nP35MmhmsB1K8" -t secp256k1
    $ okp4d keys did 02d0fe99b214aaeeb5e46ae4d65a1623a95d9d0cecd57d673f7e4c9a19ac0752bc -t secp256k1

```
okp4d keys did [pubkey] -t [{ed25519, secp256k1}] [flags]
```

### Options

```
  -h, --help          help for did
  -t, --type string   Pubkey type to decode (oneof ed25519, secp256k1) (default "secp256k1")
```

### Options inherited from parent commands

```
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) (default "test")
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --output string            Output format (text|json) (default "text")
```

### SEE ALSO

* [okp4d keys](okp4d_keys.md)	 - Manage your application's keys
