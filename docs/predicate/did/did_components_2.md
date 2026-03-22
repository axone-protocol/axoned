---
sidebar_position: 1
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# did_components/2

## Module

This predicate is provided by `did.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/did.pl').
```

## Description

Also supports the reverse mode:

```prolog
did_components(-DID:atom, +Parsed) is det.
```

Parse or reconstruct a DID / DID URL compliant with W3C DID Core.

```prolog
Parsed = did(Method, MethodSpecificId, Path, Query, Fragment)
```

where:

- Method is an atom.
- MethodSpecificId is an atom.
- Path is an atom including its leading `/` when present, otherwise left unbound.
- Query is a raw atom without its leading `?` when present, otherwise left unbound.
- Fragment is a raw atom without its leading `#` when present, otherwise left unbound.

Components are preserved raw. No percent-decoding or URI normalization is performed.

## Signature

```text
did_components(+DID:atom, -Parsed) is det
```

## Examples

### Parse a DID URL into raw DID components

This scenario demonstrates how to decompose a DID URL into a `did/5` structured term.
Path is preserved with its leading `/`, while query and fragment are preserved raw without their leading separators.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/did.pl'),
did_components(
  'did:example:123456/path?versionId=1#auth-key',
  did(Method, MethodSpecificId, Path, Query, Fragment)
).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 11932
answer:
  has_more: false
  variables: ["Method", "MethodSpecificId", "Path", "Query", "Fragment"]
  results:
  - substitutions:
    - variable: Method
      expression: "example"
    - variable: MethodSpecificId
      expression: "'123456'"
    - variable: Path
      expression: "'/path'"
    - variable: Query
      expression: "'versionId=1'"
    - variable: Fragment
      expression: "'auth-key'"
```

### Reconstruct a DID URL from raw DID components

This scenario demonstrates the reverse mode of did_components/2, reconstructing a DID URL from a `did/5` term.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/did.pl'),
did_components(DID, did(example, '123456', '/foo/bar', 'versionId=1', test)).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 10814
answer:
  has_more: false
  variables: ["DID"]
  results:
  - substitutions:
    - variable: DID
      expression: "'did:example:123456/foo/bar?versionId=1#test'"
```

### Error on invalid DID encoding

This scenario demonstrates the error returned when the DID text does not comply with DID Core syntax.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/did.pl'),
did_components(foo, Parsed).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4315
answer:
  has_more: false
  variables: ["Parsed"]
  results:
  - error: "error(domain_error(encoding(did),foo),did_components/2)"
```

### Error on invalid raw path when reconstructing

This scenario demonstrates the error returned when a parsed DID term contains a path that is not encoded according to the selected raw representation.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/did.pl'),
did_components(DID, did(example, '123456', 'path with/space', _, _)).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 6391
answer:
  has_more: false
  variables: ["DID"]
  results:
  - error: "error(domain_error(encoding(did),path with/space),did_components/2)"
```
