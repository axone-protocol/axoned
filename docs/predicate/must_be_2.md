---
sidebar_position: 82
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# must_be/2

## Module

This predicate is provided by `error.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/error.pl').
```

## Description

Succeeds when Term satisfies Type.
Throws:

- error(instantiation_error, must_be/2) when Term is insufficiently instantiated;
- error(type_error(Type, Term), must_be/2) when Term has the wrong type;
- error(existence_error(type, Type), must_be/2) when Type is unknown.

## Signature

```text
must_be(+Type, @Term) is det
```

## Examples

### Validate an atom with must_be/2

This scenario demonstrates how to load error.pl and validate a value type with must_be/2.

Here are the steps of the scenario:

- **Given** the program:

```  prolog

```

- **Given** the query:

```  prolog
consult('/v1/lib/error.pl'),
must_be(atom, hello),
Result = ok.
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3986
answer:
  has_more: false
  variables: ["Result"]
  results:
  - substitutions:
    - variable: Result
      expression: ok
```

### must_be/2 throws instantiation_error for unbound values

This scenario demonstrates that must_be/2 raises an instantiation error when the checked value is a variable.

Here are the steps of the scenario:

- **Given** the program:

```  prolog

```

- **Given** the query:

```  prolog
consult('/v1/lib/error.pl'),
must_be(atom, X).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3997
answer:
  has_more: false
  variables: ["X"]
  results:
  - error: "error(instantiation_error,must_be/2)"
```
