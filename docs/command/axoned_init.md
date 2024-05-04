## axoned init

Initialize private validator, p2p, genesis, and application configuration files

### Synopsis

Initialize validators's and node's configuration files.

```
axoned init [moniker] [flags]
```

### Options

```
      --chain-id string        genesis file chain-id, if left blank will be randomly created (default "axoned")
      --default-denom string   genesis file default denomination, if left blank default value is 'stake'
  -h, --help                   help for init
      --home string            node's home directory (default "/home/john/.axoned")
      --initial-height int     specify the initial block height at genesis (default 1)
  -o, --overwrite              overwrite the genesis.json file
      --recover                provide seed phrase to recover existing key instead of creating
```

### SEE ALSO

* [axoned](axoned.md)	 - Axone - Orchestration Layer for AI
