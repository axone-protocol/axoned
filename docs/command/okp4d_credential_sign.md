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
      --schema-map strings       Map original URIs to alternative URIs for resolving JSON-LD schemas. Useful for redirecting network-based URIs to local filesystem paths or other URIs. Each mapping should be in the format 'originalURI=alternativeURI'. Multiple mappings can be provided by repeating the flag. Example usage: --schema-map originalURI1=alternativeURI1 --schema-map originalURI2=alternativeURI2
```

### SEE ALSO

* [okp4d credential](okp4d_credential.md)	 - W3C Verifiable Credential
