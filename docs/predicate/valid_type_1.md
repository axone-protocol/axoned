---
sidebar_position: 132
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# valid_type/1

## Module

This predicate is provided by `type.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/type.pl').
```

## Description

Succeeds when Type is a valid, concrete type specification.
Recursively validates parameterized types (e.g., list(integer), list(list(byte))).
Fails when Type is unbound or contains unbound type parameters.

## Signature

```text
valid_type(+Type) is semidet
```
