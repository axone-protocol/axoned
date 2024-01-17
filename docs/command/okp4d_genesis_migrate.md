## okp4d genesis migrate

Migrate genesis to a specified target version

### Synopsis

Migrate the source genesis into the target version and print to STDOUT

```
okp4d genesis migrate [target-version] [genesis-file] [flags]
```

### Examples

```
okp4d migrate v0.47 /path/to/genesis.json --chain-id=cosmoshub-3 --genesis-time=2019-04-22T17:00:00Z
```

### Options

```
      --chain-id string          Override chain_id with this flag (default "okp4d")
      --genesis-time string      Override genesis_time with this flag
  -h, --help                     help for migrate
      --output-document string   Exported state is written to the given file instead of STDOUT
```

### SEE ALSO

* [okp4d genesis](okp4d_genesis.md)	 - Application's genesis-related subcommands
