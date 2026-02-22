---
sidebar_position: 133
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# peek_code/2

## Module

This predicate is provided by `stdlib.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/stdlib.pl').
```

## Description

Peeks the next character code from Stream without consuming it.
Returns -1 at end of file.

## Signature

```text
peek_code(+Stream, ?Code) is det
```
