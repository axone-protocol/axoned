---
sidebar_position: 1
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# base64/2

## Module

This predicate is provided by `base64.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/base64.pl').
```

## Description

Relates a text value to its classic padded Base64 representation.

The predicate is equivalent to `base64_encoded/3` with options
`[as(atom), encoding(utf8), charset(classic), padding(true)]`.

## Signature

```text
base64(?Plain, ?Encoded) is det
```

## Examples

### Encode and decode a string into a Base64 encoded atom

This scenario demonstrates how to encode an decode a plain string into a Base64-encoded atom using the `base64/2`
predicate.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/base64.pl'),
base64('Hello world', Encoded),
base64(Decoded, 'SGVsbG8gd29ybGQ=').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 19603
answer:
  has_more: false
  variables: ["Encoded", "Decoded"]
  results:
  - substitutions:
    - variable: Encoded
      expression: "'SGVsbG8gd29ybGQ='"
    - variable: Decoded
      expression: "'Hello world'"
```
