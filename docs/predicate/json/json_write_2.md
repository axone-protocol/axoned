---
sidebar_position: 3
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# json_write/2

## Module

This predicate is provided by `json.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/json.pl').
```

## Description

Writes Term as JSON text to Stream.

## Signature

```text
json_write(+Stream, +Term) is det
```

## Examples

### Write a canonical Prolog JSON term to a stream

This scenario demonstrates writing a Prolog JSON term to the current output stream.

Here are the steps of the scenario:

- **Given** the module configuration:

```  json
{
  "limits": {
    "max_user_output_size": 50
  }
}
```

- **Given** the program:

```  prolog
:- consult('/v1/lib/json.pl').

json_write_to_output(Term) :-
  current_output(Stream),
  json_write(Stream, Term).
```

- **Given** the query:

```  prolog
json_write_to_output(json([foo=bar,ok= @(true)])).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 11668
answer:
  has_more: false
  variables:
  results:
  - substitutions:
user_output: |
  {"foo":"bar","ok":true}
```
