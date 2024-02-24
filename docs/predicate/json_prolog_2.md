---
sidebar_position: 13
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# json_prolog/2

## Description

`json_prolog/2` is a predicate that will unify a JSON string into prolog terms and vice versa.

The signature is as follows:

```text
json_prolog(?Json, ?Term) is det
```

Where:

- Json is the string representation of the json
- Term is an Atom that would be unified by the JSON representation as Prolog terms.

In addition, when passing Json and Term, this predicate return true if both result match.

## Examples

```text
# JSON conversion to Prolog.
- json_prolog('{"foo": "bar"}', json([foo-bar])).
```
