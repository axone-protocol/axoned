## okp4d config migrate

Migrate Cosmos SDK app configuration file to the specified version

### Synopsis

Migrate the contents of the Cosmos SDK app configuration (app.toml) to the specified version.
The output is written in-place unless --stdout is provided.
In case of any error in updating the file, no output is written.

```
okp4d config migrate [target-version] <app-toml-path> (options) [flags]
```

### Options

```
  -h, --help            help for migrate
      --skip-validate   skip configuration validation (allows to migrate unknown configurations)
      --stdout          print the updated config to stdout
      --verbose         log changes to stderr
```

### SEE ALSO

* [okp4d config](okp4d_config.md)	 - Utilities for managing application configuration
