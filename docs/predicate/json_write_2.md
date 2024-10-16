---
sidebar_position: 19
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# json_write/2

## Description

`json_write/2` is a predicate that writes a Prolog term as a JSON to a stream.

The JSON object is of the same format as produced by json\_read/2.

The signature is as follows:

```text
json_write(+Stream, +Term) is det
```

Where:

- Stream is the output stream to which the JSON is written.
- Term is the Prolog term that represents the JSON structure.
