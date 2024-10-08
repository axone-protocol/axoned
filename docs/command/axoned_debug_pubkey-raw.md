## axoned debug pubkey-raw

Decode a ED25519 or secp256k1 pubkey from hex, base64, or bech32

### Synopsis

Decode a pubkey from hex, base64, or bech32.

```
axoned debug pubkey-raw [pubkey] -t [{ed25519, secp256k1}] [flags]
```

### Examples

```

axoned debug pubkey-raw 8FCA9D6D1F80947FD5E9A05309259746F5F72541121766D5F921339DD061174A
axoned debug pubkey-raw j8qdbR+AlH/V6aBTCSWXRvX3JUESF2bV+SEzndBhF0o=
axoned debug pubkey-raw cosmospub1zcjduepq3l9f6mglsz28l40f5pfsjfvhgm6lwf2pzgtkd40eyyeem5rpza9q47axrz
			
```

### Options

```
  -h, --help          help for pubkey-raw
  -t, --type string   Pubkey type to decode (oneof secp256k1, ed25519) (default "ed25519")
```

### SEE ALSO

* [axoned debug](axoned_debug.md)	 - Tool for helping with debugging your application
