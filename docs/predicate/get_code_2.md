---
sidebar_position: 118
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# get_code/2

## Module

This predicate is provided by `stdlib.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/stdlib.pl').
```

## Description

Reads the next character code from Stream.
Returns -1 at end of file.

## Signature

```text
get_code(+Stream, ?Code) is det
```
