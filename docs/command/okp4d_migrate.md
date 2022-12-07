## okp4d migrate

Migrate genesis to a specified target version

### Synopsis

Migrate the source genesis into the target version and print to STDOUT.

Example:
$ okp4d migrate v0.36 /path/to/genesis.json --chain-id=cosmoshub-3 --genesis-time=2019-04-22T17:00:00Z


```
okp4d migrate [target-version] [genesis-file] [flags]
```

### Options

```
      --chain-id string       override chain_id with this flag (default "okp4d")
      --genesis-time string   override genesis_time with this flag
  -h, --help                  help for migrate
```

### SEE ALSO

* [okp4d](okp4d.md)	 - OKP4 Daemon 👹

