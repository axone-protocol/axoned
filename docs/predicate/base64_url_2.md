---
sidebar_position: 11
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# base64_url/2

## Description

`base64_url/2` is a predicate that unifies a string to a Base64 URL\-safe encoded string.

Encoded values are safe for use in URLs and filenames: "\+" is replaced by "\-", "/" by "\_", and padding is omitted.

The signature is as follows:

```text
base64url(+Plain, -Encoded) is det
base64url(-Plain, +Encoded) is det
```

Where:

- Plain is an atom, a list of characters, or character codes representing the unencoded text.
- Encoded is an atom, a list of characters, or character codes representing the Base64 URL\-safe encoded form.

The predicate is equivalent to base64\_encoded/3 with options: \[as\(atom\), encoding\(utf8\), charset\(url\), padding\(false\)\].
