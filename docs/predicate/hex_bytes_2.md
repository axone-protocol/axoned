---
sidebar_position: 67
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# hex_bytes/2

## Module

This predicate is provided by `crypto.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/crypto.pl').
```

## Description

Relates a hexadecimal text representation to a list of bytes.

- Hex may be an atom, a list of characters, or a list of character codes.
- Bytes is a proper list of integers in [0,255].
- At least one argument must be instantiated.
- When converting Bytes to Hex, Hex is returned as a lowercase atom.

## Signature

```text
hex_bytes(?Hex, ?Bytes) is det
```

## Examples

### Decode a hexadecimal atom into bytes

This scenario demonstrates how to decode a hexadecimal atom into a list of bytes.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/crypto.pl'),
hex_bytes('501ACE', Bytes).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 5016
answer:
  has_more: false
  variables: ["Bytes"]
  results:
  - substitutions:
    - variable: Bytes
      expression: "[80,26,206]"
```

### Encode bytes into a hexadecimal atom

This scenario demonstrates how to encode a list of bytes into a lowercase hexadecimal atom.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/crypto.pl'),
hex_bytes(Hex, [80,26,206]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 8173
answer:
  has_more: false
  variables: ["Hex"]
  results:
  - substitutions:
    - variable: Hex
      expression: "'501ace'"
```
