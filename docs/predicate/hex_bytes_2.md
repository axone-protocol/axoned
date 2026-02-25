---
sidebar_position: 65
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# hex_bytes/2

## Description

`hex_bytes/2` is a predicate that unifies hexadecimal encoded bytes to a list of bytes.

The signature is as follows:

```text
hex_bytes(?Hex, ?Bytes) is det
```

Where:

- Hex is an Atom, string or list of characters in hexadecimal encoding.
- Bytes is the list of numbers between 0 and 255 that represent the sequence of bytes.

## Examples

```text
# Convert hexadecimal atom to list of bytes.
- hex_bytes('2c26b46b68ffc68ff99b453c1d3041341342d706483bfa0f98a5e886266e7ae', Bytes).
```
