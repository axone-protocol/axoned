---
sidebar_position: 116
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# foldl/5

## Module

This predicate is provided by `apply.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/apply.pl').
```

## Description

Left-folds two lists in lockstep using Goal.
Goal is called as call(Goal, Elem1, Elem2, Acc0, Acc1).

## Signature

```text
foldl(:Goal, +List1, +List2, +V0, -V) is det
```
