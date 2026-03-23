---
sidebar_position: 1
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# uri_encoded/3

## Module

This predicate is provided by `uri.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/uri.pl').
```

## Description

uri_encoded(+Component, -Value, +Encoded) is det.

Encoded is the URI encoding for Value.

Component specifies the URI component where the value is used. It is one of
`query_value`, `fragment`, `path` or `segment`.

Value and Encoded may be atoms, lists of characters, or lists of character
codes. Generated values are returned as atoms.

## Signature

```text
uri_encoded(+Component, +Value, -Encoded) is det
```

## Examples

### Decode a raw path atom

This scenario demonstrates how to decode a raw URI path into plain text.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/uri.pl'),
uri_encoded(path, Decoded, foo).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4709
answer:
  has_more: false
  variables: ["Decoded"]
  results:
  - substitutions:
    - variable: Decoded
      expression: "foo"
```

### Encode a query value with a space

This scenario demonstrates how to percent-encode a query value.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/uri.pl'),
uri_encoded(query_value, 'foo bar', Encoded).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 5716
answer:
  has_more: false
  variables: ["Encoded"]
  results:
  - substitutions:
    - variable: Encoded
      expression: "'foo%20bar'"
```
