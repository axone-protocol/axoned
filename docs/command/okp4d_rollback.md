## okp4d rollback

rollback cosmos-sdk and tendermint state by one height

### Synopsis


A state rollback is performed to recover from an incorrect application state transition,
when Tendermint has persisted an incorrect app hash and is thus unable to make
progress. Rollback overwrites a state at height n with the state at height n - 1.
The application also rolls back to height n - 1. No blocks are removed, so upon
restarting Tendermint the transactions in block n will be re-executed against the
application.


```
okp4d rollback [flags]
```

### Options

```
  -h, --help          help for rollback
      --home string   The application home directory (default "/Users/chris/.okp4d")
```

### SEE ALSO

* [okp4d](okp4d.md)	 - OKP4 Daemon 👹

