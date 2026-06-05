---
sidebar_position: 1
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# json_prolog/2

## Module

This predicate is provided by `json.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/json.pl').
```

## Description

Relates JSON text with its canonical Prolog representation.

Json is text: an atom, a list of characters, or a list of character codes.

The canonical representation for Term is:

- JSON objects are represented as `json(NameValueList)`;
- JSON arrays are represented as Prolog lists;
- JSON strings are represented as atoms;
- JSON numbers are represented as numbers;
- JSON booleans and null are represented as `@(true)`, `@(false)`, and `@(null)`.

## Signature

```text
json_prolog(?Json, ?Term) is det
```

## Examples

### Decode JSON text into a canonical Prolog term

This scenario demonstrates how JSON objects, strings, and booleans are represented in Prolog.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/json.pl'),
json_prolog('{"foo":"bar","ok":true}', Term).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 12069
answer:
  has_more: false
  variables: ["Term"]
  results:
  - substitutions:
    - variable: Term
      expression: "json([foo=bar,ok= @(true)])"
```

### Encode a canonical Prolog term as JSON text

This scenario demonstrates how a canonical Prolog JSON object is encoded as compact JSON text.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/json.pl'),
json_prolog(Json, json([foo=bar,ok= @(true)])).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 11049
answer:
  has_more: false
  variables: ["Json"]
  results:
  - substitutions:
    - variable: Json
      expression: "'{\"foo\":\"bar\",\"ok\":true}'"
```
