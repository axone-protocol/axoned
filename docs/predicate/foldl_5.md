---
sidebar_position: 112
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

## Examples

### Fold two lists in lockstep to compute dot product

This scenario demonstrates how to use foldl/5 to fold two lists simultaneously.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
add_product(X, Y, Acc0, Acc) :- Acc is Acc0 + (X * Y).
```

- **Given** the query:

```  prolog
consult('/v1/lib/apply.pl'),
foldl(add_product, [1,2,3], [4,5,6], 0, DotProduct).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3988
answer:
  has_more: false
  variables: ["DotProduct"]
  results:
  - substitutions:
    - variable: DotProduct
      expression: 32
```

### Fold two empty lists returns the initial accumulator

Here are the steps of the scenario:

- **Given** the program:

```  prolog
add_product(X, Y, Acc0, Acc) :- Acc is Acc0 + (X * Y).
```

- **Given** the query:

```  prolog
consult('/v1/lib/apply.pl'),
foldl(add_product, [], [], 99, Result).
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
      expression: 99
```

### Fold two lists to build a pair list

Here are the steps of the scenario:

- **Given** the program:

```  prolog
make_pair(X, Y, Acc0, [[X,Y]|Acc0]).
```

- **Given** the query:

```  prolog
consult('/v1/lib/apply.pl'),
foldl(make_pair, [a,b], [1,2], [], Pairs).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3982
answer:
  has_more: false
  variables: ["Pairs"]
  results:
  - substitutions:
    - variable: Pairs
      expression: "[[b,2],[a,1]]"
```
