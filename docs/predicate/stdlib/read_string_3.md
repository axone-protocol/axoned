---
sidebar_position: 68
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# read_string/3

## Module

Built-in predicate.

## Description

Reads characters from Stream and unifies String with an atom containing the
text read. Length is unified with the number of UTF-8 bytes read. When Length
is instantiated to a positive integer, reading stops once at least that many
bytes have been read.

## Signature

```text
read_string(+Stream, ?Length, -String) is det
```

## Examples

### Read a text stream into an atom and byte length

This scenario demonstrates reading all text from a stream while counting UTF-8 bytes.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
read_from_echo(Text, Length, String) :-
  open('/v1/dev/echo', read_write, Stream, [type(text)]),
  write(Stream, Text),
  read_string(Stream, Length, String),
  close(Stream).
```

- **Given** the query:

```  prolog
read_from_echo('aé', Length, String).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 5505
answer:
  has_more: false
  variables: ["Length", "String"]
  results:
  - substitutions:
    - variable: Length
      expression: 3
    - variable: String
      expression: "aé"
```
