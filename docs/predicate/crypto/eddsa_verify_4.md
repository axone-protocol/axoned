---
sidebar_position: 3
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# eddsa_verify/4

## Module

This predicate is provided by `crypto.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/crypto.pl').
```

## Description

Succeeds when Signature is a valid EdDSA signature for Data and PubKey.

PubKey and Signature are lists of bytes.

Options may be a single option term or a list of option terms:

- `type(+Algorithm)` selects `ed25519` (default);
- `encoding(+Encoding)` selects how Data is interpreted, defaulting to `hex`.

Supported encodings are `hex`, `octet`, `utf8`, and `text`.

## Signature

```text
eddsa_verify(+PubKey, +Data, +Signature, +Options) is semidet
```

## Examples

### Verify an Ed25519 signature with default options

This scenario demonstrates how to verify an Ed25519 signature over hexadecimal data.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
valid_ed25519(Verified) :-
  hex_bytes('53167ac3fc4b720daa45b04fc73fe752578fa23a10048422d6904b7f4f7bba5a', PubKey),
  hex_bytes('889bcfd331e8e43b5ebf430301dffb6ac9e2fce69f6227b43552fe3dc8cc1ee00c1cc53452a8712e9d5f80086dff8cf4999c1b93ed6c6e403c09334cb61ddd0b', Signature),
  eddsa_verify(PubKey, '9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', Signature, []),
  Verified = true.
```

- **Given** the query:

```  prolog
consult('/v1/lib/crypto.pl'),
valid_ed25519(Verified).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 235381
answer:
  has_more: false
  variables: ["Verified"]
  results:
  - substitutions:
    - variable: Verified
      expression: true
```

### Reject an invalid Ed25519 signature

This scenario demonstrates that eddsa_verify/4 fails when the signature does not match the data.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
invalid_ed25519 :-
  hex_bytes('53167ac3fc4b720daa45b04fc73fe752578fa23a10048422d6904b7f4f7bba5a', PubKey),
  hex_bytes('889bcfd331e8e43b5ebf430301dffb6ac9e2fce69f6227b43552fe3dc8cc1ee00c1cc53452a8712e9d5f80086dff8cf4999c1b93ed6c6e403c09334cb61ddd0b', Signature),
  eddsa_verify(PubKey, '9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9e', Signature, []).
```

- **Given** the query:

```  prolog
consult('/v1/lib/crypto.pl'),
invalid_ed25519.
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 235349
answer:
  has_more: false
  variables:
  results:
```

### Reject an unsupported EdDSA algorithm

This scenario demonstrates that eddsa_verify/4 rejects algorithms outside the EdDSA family.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
unsupported_eddsa :-
  eddsa_verify([], '', [], [type(secp256k1)]).
```

- **Given** the query:

```  prolog
consult('/v1/lib/crypto.pl'),
unsupported_eddsa.
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4809
answer:
  has_more: false
  variables:
  results:
  - error: "error(type_error(cryptographic_algorithm,secp256k1),eddsa_verify/4)"
```
