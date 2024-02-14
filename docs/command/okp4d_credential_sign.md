## okp4d credential sign

Sign a W3C Verifiable Credential provided as a file or stdin

### Synopsis

Sign a W3C Verifiable Credential;

It will read a verifiable credential from a file (or stdin), sign it, and print the JSON-LD signed credential to stdout.

```
okp4d credential sign [file] [flags]
```

### Options

```
      --date string              Date of the signature provided in RFC3339 format. If not provided, current time will be used
      --from string              Name or address of private key with which to sign
  -h, --help                     help for sign
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) (default "test")
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --overwrite                Overwrite existing signatures with a new one. If disabled, new signature will be appended
```

### SEE ALSO

* [okp4d credential](okp4d_credential.md)	 - W3C Verifiable Credential
