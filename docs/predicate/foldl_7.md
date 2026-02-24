---
sidebar_position: 114
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# foldl/7

## Module

This predicate is provided by `apply.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/apply.pl').
```

## Description

Left-folds four lists in lockstep using Goal.
Goal is called as call(Goal, Elem1, Elem2, Elem3, Elem4, Acc0, Acc1).

## Signature

```text
foldl(:Goal, +List1, +List2, +List3, +List4, +V0, -V) is det
```

## Examples

### Fold four lists in lockstep to compute a combined result

This scenario demonstrates how to use foldl/7 to fold four lists simultaneously.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
quad_sum(W, X, Y, Z, Acc0, Acc) :- Acc is Acc0 + W + X + Y + Z.
```

- **Given** the query:

```  prolog
consult('/v1/lib/apply.pl'),
foldl(quad_sum, [1,2], [3,4], [5,6], [7,8], 0, Result).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3984
answer:
  has_more: false
  variables: ["Result"]
  results:
  - substitutions:
    - variable: Result
      expression: 36
```

### Fold four empty lists returns the initial accumulator

Here are the steps of the scenario:

- **Given** the program:

```  prolog
quad_sum(W, X, Y, Z, Acc0, Acc) :- Acc is Acc0 + W + X + Y + Z.
```

- **Given** the query:

```  prolog
consult('/v1/lib/apply.pl'),
foldl(quad_sum, [], [], [], [], 100, Result).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3976
answer:
  has_more: false
  variables: ["Result"]
  results:
  - substitutions:
    - variable: Result
      expression: 100
```

### Fold four lists to build a structured result

Here are the steps of the scenario:

- **Given** the program:

```  prolog
make_quad(W, X, Y, Z, Acc0, [[W,X,Y,Z]|Acc0]).
```

- **Given** the query:

```  prolog
consult('/v1/lib/apply.pl'),
foldl(make_quad, [a], [1], [x], [true], [], Quads).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3979
answer:
  has_more: false
  variables: ["Quads"]
  results:
  - substitutions:
    - variable: Quads
      expression: "[[a,1,x,true]]"
```
