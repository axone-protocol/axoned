---
sidebar_position: 2
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# ecdsa_verify/4

## Module

This predicate is provided by `crypto.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/crypto.pl').
```

## Description

Succeeds when Signature is a valid ECDSA signature for Data and PubKey.

PubKey is the compressed public key as a list of bytes. Signature is the ASN.1
encoded signature as a list of bytes.

Options may be a single option term or a list of option terms:

- `type(+Algorithm)` selects `secp256r1` (default) or `secp256k1`;
- `encoding(+Encoding)` selects how Data is interpreted, defaulting to `hex`.

Supported encodings are `hex`, `octet`, `utf8`, and `text`.

## Signature

```text
ecdsa_verify(+PubKey, +Data, +Signature, +Options) is semidet
```

## Examples

### Verify a secp256r1 signature with default options

This scenario demonstrates how to verify an ECDSA signature over hexadecimal data.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
valid_secp256r1(Verified) :-
  hex_bytes('0213c8426be471e55506f7ce4f7df557a42e310df09f92eb732ca3085e797cef9b', PubKey),
  hex_bytes('30450220099e6f9dd218e0e304efa7a4224b0058a8e3aec73367ec239bee4ed8ed7d85db022100b504d3d0d2e879b04705c0e5a2b40b0521a5ab647ea207bd81134e1a4eb79e47', Signature),
  ecdsa_verify(PubKey, 'e50c26e89f734b2ee12041ff27874c901891f74a0f0cf470333312a3034ce3be', Signature, []),
  Verified = true.
```

- **Given** the query:

```  prolog
consult('/v1/lib/crypto.pl'),
valid_secp256r1(Verified).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 246079
answer:
  has_more: false
  variables: ["Verified"]
  results:
  - substitutions:
    - variable: Verified
      expression: true
```

### Verify a secp256k1 signature

This scenario demonstrates how to select secp256k1 explicitly.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
valid_secp256k1(Verified) :-
  hex_bytes('026b5450187ee9c63ba9e42cb6018d8469c903aca116178e223de76e49fe63b71c', PubKey),
  hex_bytes('304402201448201bb4408549b0997f4b9ad9ed36f3cf8bb9c433fc7f3ba48c6b6e39476e022053f7d056f7ffeab9a79f3a36bc2ba969ddd530a3a1495d1ed7bba00039820223', Signature),
  ecdsa_verify(PubKey, 'dece063885d3648078f903b6a3e8989f649dc3368cd9c8d69755ed9dcb6a0995', Signature, [type(secp256k1)]),
  Verified = true.
```

- **Given** the query:

```  prolog
consult('/v1/lib/crypto.pl'),
valid_secp256k1(Verified).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 250523
answer:
  has_more: false
  variables: ["Verified"]
  results:
  - substitutions:
    - variable: Verified
      expression: true
```

### Verify ECDSA signatures with explicit algorithms

This scenario demonstrates how to select both supported ECDSA algorithms explicitly.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
valid_explicit_secp256r1 :-
  hex_bytes('0213c8426be471e55506f7ce4f7df557a42e310df09f92eb732ca3085e797cef9b', PubKey),
  hex_bytes('30450220099e6f9dd218e0e304efa7a4224b0058a8e3aec73367ec239bee4ed8ed7d85db022100b504d3d0d2e879b04705c0e5a2b40b0521a5ab647ea207bd81134e1a4eb79e47', Signature),
  ecdsa_verify(PubKey, 'e50c26e89f734b2ee12041ff27874c901891f74a0f0cf470333312a3034ce3be', Signature, [type(secp256r1)]).

valid_explicit_secp256k1 :-
  hex_bytes('026b5450187ee9c63ba9e42cb6018d8469c903aca116178e223de76e49fe63b71c', PubKey),
  hex_bytes('304402201448201bb4408549b0997f4b9ad9ed36f3cf8bb9c433fc7f3ba48c6b6e39476e022053f7d056f7ffeab9a79f3a36bc2ba969ddd530a3a1495d1ed7bba00039820223', Signature),
  ecdsa_verify(PubKey, 'dece063885d3648078f903b6a3e8989f649dc3368cd9c8d69755ed9dcb6a0995', Signature, [type(secp256k1)]).
```

- **Given** the query:

```  prolog
consult('/v1/lib/crypto.pl'),
valid_explicit_secp256r1,
valid_explicit_secp256k1,
Verified = true.
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 493294
answer:
  has_more: false
  variables: ["Verified"]
  results:
  - substitutions:
    - variable: Verified
      expression: true
```

### Reject an invalid ECDSA signature

This scenario demonstrates that ecdsa_verify/4 fails when the signature does not match the data.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
invalid_secp256r1 :-
  hex_bytes('0213c8426be471e55506f7ce4f7df557a42e310df09f92eb732ca3085e797cef9b', PubKey),
  hex_bytes('30450220099e6f9dd218e0e304efa7a4224b0058a8e3aec73367ec239bee4ed8ed7d85db022100b504d3d0d2e879b04705c0e5a2b40b0521a5ab647ea207bd81134e1a4eb79e48', Signature),
  ecdsa_verify(PubKey, 'e50c26e89f734b2ee12041ff27874c901891f74a0f0cf470333312a3034ce3be', Signature, []).
```

- **Given** the query:

```  prolog
consult('/v1/lib/crypto.pl'),
invalid_secp256r1.
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 246047
answer:
  has_more: false
  variables:
  results:
```

### Reject an unsupported ECDSA algorithm

This scenario demonstrates that ecdsa_verify/4 rejects algorithms outside the ECDSA family.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
unsupported_ecdsa :-
  ecdsa_verify([], '', [], [type(ed25519)]).
```

- **Given** the query:

```  prolog
consult('/v1/lib/crypto.pl'),
unsupported_ecdsa.
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4841
answer:
  has_more: false
  variables:
  results:
  - error: "error(type_error(cryptographic_algorithm,ed25519),ecdsa_verify/4)"
```
