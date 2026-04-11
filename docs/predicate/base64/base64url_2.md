---
sidebar_position: 2
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# base64url/2

## Module

This predicate is provided by `base64.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/base64.pl').
```

## Description

base64url(-Plain, +Encoded) is det.

Relates a text value to its URL-safe Base64 representation.

The predicate is equivalent to `base64_encoded/3` with options
`[as(atom), encoding(utf8), charset(url), padding(false)]`.

## Signature

```text
base64url(+Plain, -Encoded) is det
```

## Examples

### Encode and decode a string into a Base64 encoded atom in URL-Safe mode

This scenario demonstrates how to encode an decode a plain string into a Base64-encoded atom using the `base64url/2`
predicate.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/base64.pl'),
base64url('<<???>>', Encoded),
base64url(Decoded, 'PDw_Pz8-Pg').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 15447
answer:
  has_more: false
  variables: ["Encoded", "Decoded"]
  results:
  - substitutions:
    - variable: Encoded
      expression: "'PDw_Pz8-Pg'"
    - variable: Decoded
      expression: "<<???>>"
```
