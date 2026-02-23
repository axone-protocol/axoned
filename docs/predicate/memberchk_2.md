---
sidebar_position: 133
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# memberchk/2

## Module

Built-in predicate.

## Description

Succeeds if Elem is a member of List. This is a deterministic predicate
that commits to the first unification and does not leave a choice point.
Useful when List is ground and you only need to check membership once.

## Signature

```text
memberchk(?Elem, +List) is semidet
```
