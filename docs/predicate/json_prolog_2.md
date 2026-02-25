---
sidebar_position: 67
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# json_prolog/2

## Description

`json_prolog/2` is a predicate that unifies a JSON into a prolog term and vice versa.

The signature is as follows:

```text
json_prolog(?Json, ?Term) is det
```

Where:

- Json is the textual representation of the JSON, as either an atom, a list of character codes, or a list of characters.
- Term is the Prolog term that represents the JSON structure.

## JSON canonical representation

The canonical representation for Term is:

- A JSON object is mapped to a Prolog term json\(NameValueList\), where NameValueList is a list of Name=Value key values. Name is an atom created from the JSON string.
- A JSON array is mapped to a Prolog list of JSON values.
- A JSON string is mapped to a Prolog atom.
- A JSON number is mapped to a Prolog number.
- The JSON constants true and false are mapped to @\(true\) and @\(false\).
- The JSON constant null is mapped to the Prolog term @\(null\).

## Examples

```text
# JSON conversion to Prolog.
- json_prolog('{"foo": "bar"}', json([foo=bar])).
```
