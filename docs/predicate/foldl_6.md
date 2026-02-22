---
sidebar_position: 117
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# foldl/6

## Module

This predicate is provided by `apply.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/apply.pl').
```

## Description

Left-folds three lists in lockstep using Goal.
Goal is called as call(Goal, Elem1, Elem2, Elem3, Acc0, Acc1).

## Signature

```text
foldl(:Goal, +List1, +List2, +List3, +V0, -V) is det
```
