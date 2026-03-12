---
sidebar_position: 130
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# with_context/2

## Module

This predicate is provided by `error.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/error.pl').
```

## Description

Executes Goal and, if it throws a Prolog error term, rethrows the same formal
error with Context as the error context.

Non-error exceptions are rethrown unchanged.

## Signature

```text
with_context(+Context, :Goal) is det
```

## Examples

### with_context/2 rewrites the error context of the wrapped goal

This scenario demonstrates how with_context/2 preserves the formal error while replacing its context.

Here are the steps of the scenario:

- **Given** the program:

```  prolog

```

- **Given** the query:

```  prolog
consult('/v1/lib/error.pl'),
with_context(example/1, must_be(atom, 42)).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4379
answer:
  has_more: false
  results:
  - error: "error(type_error(atom,42),example/1)"
```
