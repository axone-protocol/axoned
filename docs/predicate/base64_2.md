---
sidebar_position: 18
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# base64/2

## Description

`base64/2` is a predicate that unifies a string to a `base64/2` encoded string.

The output is returned as an atom with padding included.

The signature is as follows:

```text
base64(+Plain, -Encoded) is det
base64(-Plain, +Encoded) is det
```

Where:

- Plain is an atom, a list of characters, or character codes representing the unencoded text.
- Encoded is an atom, a list of characters, or character codes representing the `base64/2` encoded form.

The predicate is equivalent to base64\_encoded/3 with options: \[as\(atom\), encoding\(utf8\), charset\(classic\), padding\(true\)\].

## Examples

### Encode and decode a string into a Base64 encoded atom

This scenario demonstrates how to encode an decode a plain string into a Base64-encoded atom using the `base64/2`
predicate.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
base64('Hello world', Encoded),
base64(Decoded, 'SGVsbG8gd29ybGQ=').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3976
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
