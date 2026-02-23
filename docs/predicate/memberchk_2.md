---
sidebar_position: 132
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# memberchk/2

## Module

This predicate is provided by `stdlib.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/stdlib.pl').
```

## Description

Succeeds if Elem is a member of List. This is a deterministic predicate
that commits to the first unification and does not leave a choice point.
Useful when List is ground and you only need to check membership once.

## Signature

```text
memberchk(?Elem, +List) is semidet
```
