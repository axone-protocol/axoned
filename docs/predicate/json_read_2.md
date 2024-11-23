---
sidebar_position: 22
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# json_read/2

## Description

`json_read/2` is a predicate that reads a JSON from a stream and unifies it with a Prolog term.

See json\_prolog/2 for the canonical representation of the JSON term.

The signature is as follows:

```text
json_read(+Stream, ?Term) is det
```

Where:

- Stream is the input stream from which the JSON is read.
- Term is the Prolog term that represents the JSON structure.
