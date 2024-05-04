## axoned export

Export state to JSON

```
axoned export [flags]
```

### Options

```
      --for-zero-height              Export state to start at height zero (perform preproccessing)
      --height int                   Export state from a particular height (-1 means latest height) (default -1)
  -h, --help                         help for export
      --home string                  The application home directory (default "/home/john/.axoned")
      --jail-allowed-addrs strings   Comma-separated list of operator addresses of jailed validators to unjail
      --modules-to-export strings    Comma-separated list of modules to export. If empty, will export all modules
      --output-document string       Exported state is written to the given file instead of STDOUT
```

### SEE ALSO

* [axoned](axoned.md)	 - Axone - Orchestration Layer for AI
