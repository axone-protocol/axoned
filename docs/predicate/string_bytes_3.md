---
sidebar_position: 93
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# string_bytes/3

## Description

`string_bytes/3` is a predicate that unifies a string with a list of bytes, returning true when the \(Unicode\) String is represented by Bytes in Encoding.

The signature is as follows:

```text
string_bytes(?String, ?Bytes, +Encoding)
```

Where:

- String is the string to convert to bytes. It can be an Atom, string or list of characters codes.
- Bytes is the list of numbers between 0 and 255 that represent the sequence of bytes.
- Encoding is the encoding to use for the conversion.

Encoding can be one of the following: \- 'text' considers the string as a sequence of Unicode characters. \- 'octet' considers the string as a sequence of bytes. \- '\<encoding\>' considers the string as a sequence of characters in the given encoding.

At least one of String or Bytes must be instantiated.

## Examples

```text
# Convert a string to a list of bytes.
- string_bytes('Hello World', Bytes, octet).

# Convert a list of bytes to a string.
- string_bytes(String, [72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100], octet).
```
