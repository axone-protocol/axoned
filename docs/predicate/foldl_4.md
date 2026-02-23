---
sidebar_position: 116
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# foldl/4

## Module

This predicate is provided by `apply.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/apply.pl').
```

## Description

Left-folds List using Goal.
Goal is called as call(Goal, Elem, Acc0, Acc1).

## Signature

```text
foldl(:Goal, +List, +V0, -V) is det
```

## Examples

### Fold a list of integers into a sum

This scenario demonstrates how to load apply.pl and use foldl/4 to aggregate a list with an accumulator.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
sum(Elem, Acc0, Acc) :- Acc is Acc0 + Elem.
```

- **Given** the query:

```  prolog
consult('/v1/lib/apply.pl'),
foldl(sum, [1,2,3,4], 0, Total).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3992
answer:
  has_more: false
  variables: ["Total"]
  results:
  - substitutions:
    - variable: Total
      expression: 10
```

### Fold an empty list returns the initial accumulator

Here are the steps of the scenario:

- **Given** the program:

```  prolog
sum(Elem, Acc0, Acc) :- Acc is Acc0 + Elem.
```

- **Given** the query:

```  prolog
consult('/v1/lib/apply.pl'),
foldl(sum, [], 42, Total).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3976
answer:
  has_more: false
  variables: ["Total"]
  results:
  - substitutions:
    - variable: Total
      expression: 42
```

### Fold with non-numeric accumulator (list building)

Here are the steps of the scenario:

- **Given** the program:

```  prolog
cons(Elem, Acc0, [Elem|Acc0]).
```

- **Given** the query:

```  prolog
consult('/v1/lib/apply.pl'),
foldl(cons, [a,b,c], [], Result).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3985
answer:
  has_more: false
  variables: ["Result"]
  results:
  - substitutions:
    - variable: Result
      expression: "[c,b,a]"
```
