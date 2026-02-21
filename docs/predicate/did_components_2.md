---
sidebar_position: 16
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# did_components/2

## Description

`did_components/2` is a predicate which breaks down a DID into its components according to the [W3C DID](<https://w3c.github.io/did-core>) specification.

The signature is as follows:

```text
did_components(+DID, -Components) is det
did_components(-DID, +Components) is det
```

where:

- DID represent DID URI, given as an Atom, compliant with [W3C DID](<https://w3c.github.io/did-core>) specification.
- Components is a compound Term in the format did\(Method, ID, Path, Query, Fragment\), aligned with the [DID syntax](<https://w3c.github.io/did-core/#did-syntax>), where: Method is the method name, ID is the method\-specific identifier, Path is the path component, Query is the query component and Fragment is the fragment component. Values are given as an Atom and are url encoded. For any component not present, its value will be null and thus will be left as an uninstantiated variable.

## Examples

```text
# Decompose a DID into its components.
- did_components('did:example:123456?versionId=1', did_components(Method, ID, Path, Query, Fragment)).

# Reconstruct a DID from its components.
- did_components(DID, did_components('example', '123456', _, 'versionId=1', _42)).
```
