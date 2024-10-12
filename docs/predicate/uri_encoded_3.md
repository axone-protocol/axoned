---
sidebar_position: 24
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# uri_encoded/3

## Description

`uri_encoded/3` is a predicate that unifies the given URI component with the given encoded or decoded string.

The signature is as follows:

```text
uri_encoded(+Component, +Value, -Encoded) is det
uri_encoded(+Component, -Value, +Encoded) is det
```

Where:

- Component represents the component of the URI to be escaped. It can be the atom 'query\_path', 'fragment', 'path' or 'segment'.
- Decoded represents the decoded string to be escaped.
- Encoded represents the encoded string.

For more information on URI encoding, refer to [RFC 3986](<https://datatracker.ietf.org/doc/html/rfc3986#section-2.1>).

## Examples

```text
# Escape the given string to be used in the path component.
- uri_encoded(path, "foo/bar", Encoded).

# Unescape the given string to be used in the path component.
- uri_encoded(path, Decoded, foo%2Fbar).
```
