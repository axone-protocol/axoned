---
sidebar_position: 1
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# crypto_data_hash/3

## Module

This predicate is provided by `crypto.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/crypto.pl').
```

## Description

Computes the Hash of Data using a configured hashing algorithm.

Options may be a single option term or a list of option terms:

- `algorithm(+Algorithm)` selects `sha256` (default), `sha512`, or `md5`;
- `encoding(+Encoding)` selects `utf8` (default), `text`, `hex`, or `octet`.

Data is interpreted according to Encoding and Hash is unified with the
resulting digest as a list of bytes.

## Signature

```text
crypto_data_hash(+Data, ?Hash, +Options) is det
```

## Examples

### Compute a SHA-256 hash with default options

This scenario demonstrates how to compute a SHA-256 digest from text using the default options.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/crypto.pl'),
crypto_data_hash('hello world', Hash, []),
hex_bytes(Hex, Hash).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 64284
answer:
  has_more: false
  variables: ["Hash", "Hex"]
  results:
  - substitutions:
    - variable: Hash
      expression: "[185,77,39,185,147,77,62,8,165,46,82,215,218,125,171,250,196,132,239,227,122,83,128,238,144,136,247,172,226,239,205,233]"
    - variable: Hex
      expression: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
```
