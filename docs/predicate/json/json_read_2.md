---
sidebar_position: 2
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# json_read/2

## Module

This predicate is provided by `json.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/json.pl').
```

## Description

Reads JSON text from Stream and unifies Term with its canonical Prolog representation.

## Signature

```text
json_read(+Stream, ?Term) is det
```

## Examples

### Read JSON text from a stream

This scenario demonstrates reading JSON text from a text stream and decoding it into a Prolog term.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- consult('/v1/lib/json.pl').

json_read_from_echo(Json, Term) :-
  open('/v1/dev/echo', read_write, Stream, [type(text)]),
  write(Stream, Json),
  json_read(Stream, Term),
  close(Stream).
```

- **Given** the query:

```  prolog
json_read_from_echo('{"foo":"bar","items":[1,null]}', Term).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 30138
answer:
  has_more: false
  variables: ["Term"]
  results:
  - substitutions:
    - variable: Term
      expression: "json([foo=bar,items=[1.0,@(null)]])"
```
