---
sidebar_position: 65
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# has_type/2

## Module

This predicate is provided by `type.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/type.pl').
```

## Description

Succeeds when Term satisfies Type without throwing.
Fails when Type is known but Term does not match it.

## Signature

```text
has_type(+Type, @Term) is semidet
```

## Examples

### Validate a byte list with has_type/2

This scenario demonstrates how to load type.pl and check a structured type using has_type/2.

Here are the steps of the scenario:

- **Given** the program:

```  prolog

```

- **Given** the query:

```  prolog
consult('/v1/lib/type.pl'),
has_type(list(byte), [0,1,255]),
Result = ok.
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4480
answer:
  has_more: false
  variables: ["Result"]
  results:
  - substitutions:
    - variable: Result
      expression: ok
```

### has_type/2 validates integer in integer-bounded between range

This scenario demonstrates that has_type/2 succeeds for an integer within an integer-bounded range (semidet).

Here are the steps of the scenario:

- **Given** the program:

```  prolog

```

- **Given** the query:

```  prolog
consult('/v1/lib/type.pl'),
has_type(between(1, 3), 2),
Result = ok.
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4142
answer:
  has_more: false
  variables: ["Result"]
  results:
  - substitutions:
    - variable: Result
      expression: ok
```

### has_type/2 rejects float in integer-bounded between range

This scenario demonstrates that has_type/2 fails when a float is checked against an integer-bounded range.

Here are the steps of the scenario:

- **Given** the program:

```  prolog

```

- **Given** the query:

```  prolog
consult('/v1/lib/type.pl'),
has_type(between(1, 3), 2.5).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4111
answer:
  has_more: false
  results:
```

### has_type/2 validates float in float-bounded between range

This scenario demonstrates that has_type/2 succeeds when a float is within a float-bounded range.

Here are the steps of the scenario:

- **Given** the program:

```  prolog

```

- **Given** the query:

```  prolog
consult('/v1/lib/type.pl'),
has_type(between(1.0, 3.0), 2.5),
Result = ok.
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4175
answer:
  has_more: false
  variables: ["Result"]
  results:
  - substitutions:
    - variable: Result
      expression: ok
```

### has_type/2 fails when the type does not match

This scenario demonstrates that has_type/2 fails quietly when the value does not satisfy the requested type.

Here are the steps of the scenario:

- **Given** the program:

```  prolog

```

- **Given** the query:

```  prolog
consult('/v1/lib/type.pl'),
has_type(integer, hello).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4017
answer:
  has_more: false
  results:
```
