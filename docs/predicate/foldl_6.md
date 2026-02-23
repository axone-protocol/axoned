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

## Examples

### Fold three lists in lockstep to compute weighted sum

This scenario demonstrates how to use foldl/6 to fold three lists simultaneously.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
weighted_sum(X, Y, Z, Acc0, Acc) :- Acc is Acc0 + (X * Y * Z).
```

- **Given** the query:

```  prolog
consult('/v1/lib/apply.pl'),
foldl(weighted_sum, [1,2], [3,4], [5,6], 0, Result).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3992
answer:
  has_more: false
  variables: ["Result"]
  results:
  - substitutions:
    - variable: Result
      expression: 63
```
