---
sidebar_position: 79
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# string_bytes/3

## Module

Built-in predicate.

## Description

Relates text and its byte representation according to Encoding.

## Signature

```text
string_bytes(?String, ?Bytes, +Encoding) is det
```

## Examples

### Encode UTF-8 text into bytes

This scenario demonstrates converting text into its UTF-8 byte sequence.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
string_bytes('aé', Bytes, utf8).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4392
answer:
  has_more: false
  variables: ["Bytes"]
  results:
  - substitutions:
    - variable: Bytes
      expression: "[97,195,169]"
```
