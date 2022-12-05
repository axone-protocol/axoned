## okp4d keys

Manage your application's keys

### Synopsis

Keyring management commands. These keys may be in any format supported by the
Tendermint crypto library and can be used by light-clients, full nodes, or any other application
that needs to sign with a private key.

The keyring supports the following backends:

    os          Uses the operating system's default credentials store.
    file        Uses encrypted file-based keystore within the app's configuration directory.
                This keyring will request a password each time it is accessed, which may occur
                multiple times in a single command resulting in repeated password prompts.
    kwallet     Uses KDE Wallet Manager as a credentials management application.
    pass        Uses the pass command line utility to store and retrieve keys.
    test        Stores keys insecurely to disk. It does not prompt for a password to be unlocked
                and it should be use only for testing purposes.

kwallet and pass backends depend on external tools. Refer to their respective documentation for more
information:
    KWallet     https://github.com/KDE/kwallet
    pass        https://www.passwordstore.org/

The pass backend requires GnuPG: https://gnupg.org/


### Options

```
  -h, --help                     help for keys
      --home string              The application home directory (default "/home/john/.okp4d")
      --keyring-backend string   Select keyring's backend (os|file|test) (default "test")
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --output string            Output format (text|json) (default "text")
```

### SEE ALSO

* [okp4d](okp4d.md)	 - OKP4 Daemon 👹
* [okp4d keys add](okp4d_keys_add.md)	 - Add an encrypted private key (either newly generated or recovered), encrypt it, and save to <name> file
* [okp4d keys delete](okp4d_keys_delete.md)	 - Delete the given keys
* [okp4d keys export](okp4d_keys_export.md)	 - Export private keys
* [okp4d keys import](okp4d_keys_import.md)	 - Import private keys into the local keybase
* [okp4d keys list](okp4d_keys_list.md)	 - List all keys
* [okp4d keys migrate](okp4d_keys_migrate.md)	 - Migrate keys from amino to proto serialization format
* [okp4d keys mnemonic](okp4d_keys_mnemonic.md)	 - Compute the bip39 mnemonic for some input entropy
* [okp4d keys parse](okp4d_keys_parse.md)	 - Parse address from hex to bech32 and vice versa
* [okp4d keys rename](okp4d_keys_rename.md)	 - Rename an existing key
* [okp4d keys show](okp4d_keys_show.md)	 - Retrieve key information by name or address

