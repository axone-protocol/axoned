## okp4d config set

Set an application config value

### Synopsis

Set an application config value. The [config] argument must be the path of the file when using the tool standalone, otherwise it must be the name of the config file without the .toml extension.

```
okp4d config set [config] [key] [value] [flags]
```

### Options

```
  -h, --help            help for set
  -s, --skip-validate   skip configuration validation (allows to mutate unknown configurations)
      --stdout          print the updated config to stdout
  -v, --verbose         log changes to stderr
```

### SEE ALSO

* [okp4d config](okp4d_config.md)	 - Utilities for managing application configuration
