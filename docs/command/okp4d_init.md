## okp4d init

Initialize private validator, p2p, genesis, and application configuration files

### Synopsis

Initialize validators's and node's configuration files.

```
okp4d init [moniker] [flags]
```

### Options

```
      --chain-id string             genesis file chain-id, if left blank will be randomly created (default "okp4d")
  -h, --help                        help for init
      --home string                 node's home directory (default "/home/john/.okp4d")
  -o, --overwrite                   overwrite the genesis.json file
      --recover                     provide seed phrase to recover existing key instead of creating
      --staking-bond-denom string   genesis file staking bond denomination, if left blank default value is 'stake'
```

### SEE ALSO

* [okp4d](okp4d.md)	 - OKP4 Daemon 👹

