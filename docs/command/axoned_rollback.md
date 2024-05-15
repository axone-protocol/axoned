## axoned rollback

rollback Cosmos SDK and CometBFT state by one height

### Synopsis

A state rollback is performed to recover from an incorrect application state transition,
when CometBFT has persisted an incorrect app hash and is thus unable to make
progress. Rollback overwrites a state at height n with the state at height n - 1.
The application also rolls back to height n - 1. No blocks are removed, so upon
restarting CometBFT the transactions in block n will be re-executed against the
application.

```
axoned rollback [flags]
```

### Options

```
      --hard          remove last block as well as state
  -h, --help          help for rollback
      --home string   The application home directory (default "/home/john/.axoned")
```

### SEE ALSO

* [axoned](axoned.md)	 - Axone - Orchestration Layer for AI
