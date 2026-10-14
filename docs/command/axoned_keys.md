## axoned keys

Manage your application's keys

### Synopsis

Keyring management commands. These keys may be in any format supported by the
CometBFT crypto library and can be used by light-clients, full nodes, or any other application
that needs to sign with a private key.

The keyring supports the following backends:

    os          Uses the operating system's default credentials store.
    file        Uses encrypted file-based keystore within the app's configuration directory.
                This keyring will request a password each time it is accessed, which may occur
                multiple times in a single command resulting in repeated password prompts.
    kwallet     Uses KDE Wallet Manager as a credentials management application.
    pass        Uses the pass command line utility to store and retrieve keys.
    test        Stores keys insecurely to disk. It does not prompt for a password to be unlocked
                and it should be used only for testing purposes.

kwallet and pass backends depend on external tools. Refer to their respective documentation for more
information:
    KWallet     [https://github.com/KDE/kwallet](https://github.com/KDE/kwallet)
    pass        [https://www.passwordstore.org/](https://www.passwordstore.org/)

The pass backend requires GnuPG: [https://gnupg.org/](https://gnupg.org/)

### Options

```
  -h, --help                     help for keys
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) (default "test")
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --output string            Output format (text|json) (default "text")
```

### SEE ALSO

* [axoned](axoned.md)	 - Axone - Orchestration Layer for AI
* [axoned keys add](axoned_keys_add.md)	 - Add an encrypted private key (either newly generated or recovered), encrypt it, and save to &lt;name&gt; file
* [axoned keys delete](axoned_keys_delete.md)	 - Delete the given keys
* [axoned keys export](axoned_keys_export.md)	 - Export private keys
* [axoned keys import](axoned_keys_import.md)	 - Import private keys into the local keybase
* [axoned keys import-hex](axoned_keys_import-hex.md)	 - Import private keys into the local keybase
* [axoned keys list](axoned_keys_list.md)	 - List all keys
* [axoned keys list-key-types](axoned_keys_list-key-types.md)	 - List all key types
* [axoned keys migrate](axoned_keys_migrate.md)	 - Migrate keys from amino to proto serialization format
* [axoned keys mnemonic](axoned_keys_mnemonic.md)	 - Compute the bip39 mnemonic for some input entropy
* [axoned keys parse](axoned_keys_parse.md)	 - Parse address from hex to bech32 and vice versa
* [axoned keys rename](axoned_keys_rename.md)	 - Rename an existing key
* [axoned keys show](axoned_keys_show.md)	 - Retrieve key information by name or address
