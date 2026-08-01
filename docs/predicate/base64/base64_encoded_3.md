---
sidebar_position: 1
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# base64_encoded/3

## Module

This predicate is provided by `base64.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/base64.pl').
```

## Description

base64_encoded(-Plain, +Encoded, +Options) is det.

Relates an atom to its Base64-encoded atom representation as specified by
[RFC 4648](https://rfc-editor.org/rfc/rfc4648.html).

The predicate follows a functional direction:

- when `Plain` is instantiated, it encodes `Plain` into `Encoded`;
- otherwise, when `Encoded` is instantiated, it decodes `Encoded` into `Plain`;
- otherwise, it throws `instantiation_error`.

`Plain` and `Encoded` are atoms. Use explicit conversion predicates such as
`atom_chars/2`, `atom_codes/2`, or `string_bytes/3` when another textual
representation is needed.

Supported options are:

- `charset(+Charset)` where `Charset` is `classic` (default) or `url`;
- `padding(+Boolean)` where `Boolean` is `true` (default) or `false`;
- `encoding(+Encoding)` to translate between text and bytes, defaulting to `utf8`.

## Signature

```text
base64_encoded(+Plain, -Encoded, +Options) is det
```

## Examples

### Encode an atom into a Base64 encoded atom with default options

This scenario demonstrates how to encode a plain atom into its Base64 representation using the `base64_encoded/3`
predicate. The default options are used, meaning:

- The output is returned as an atom.
- Padding characters (`=`) are included (`padding(true)`).
- The classic Base64 character set is used (`charset(classic)`), not the URL-safe variant.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/base64.pl'),
base64_encoded('Hello World', X, []).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 9406
answer:
  has_more: false
  variables: ["X"]
  results:
  - substitutions:
    - variable: X
      expression: "'SGVsbG8gV29ybGQ='"
```

### Reject the removed output-shape option

This scenario demonstrates that `base64_encoded/3` has a single textual output shape: atom.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/base64.pl'),
base64_encoded('Hello World', X, [as(atom)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4407
answer:
  has_more: false
  variables: ["X"]
  results:
  - error: "error(type_error(option,as(atom)),base64_encoded/3)"
    substitutions:
```

### Encode a string into a Base64 encoded atom without padding

This scenario demonstrates how to encode a plain atom into a Base64-encoded atom using the `base64_encoded/3` predicate
with custom options. The following options are used:

- `padding(false)` – padding characters (`=`) are omitted.
- The classic Base64 character set is used by default (`charset(classic)`).

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/base64.pl'),
base64_encoded('Hello World', X, [padding(false)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 9982
answer:
  has_more: false
  variables: ["X"]
  results:
  - substitutions:
    - variable: X
      expression: "'SGVsbG8gV29ybGQ'"
```

### Encode an atom into a Base64 encoded atom in URL-Safe mode

This scenario demonstrates how to encode a plain atom into a Base64-encoded atom using the `base64_encoded/3` predicate
with URL-safe encoding. The following options are used:

- `charset(url)` – the URL-safe Base64 alphabet is used (e.g., `-` and `_` instead of `+` and `/`).
- Padding characters are included by default (`padding(true)`).

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/base64.pl'),
base64_encoded('<<???>>', Classic, [charset(classic)]),
base64_encoded('<<???>>', UrlSafe, [charset(url)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 13251
answer:
  has_more: false
  variables: ["Classic", "UrlSafe"]
  results:
  - substitutions:
    - variable: Classic
      expression: "'PDw/Pz8+Pg=='"
    - variable: UrlSafe
      expression: "'PDw_Pz8-Pg=='"
```

### Decode a Base64 encoded atom into a plain atom

This scenario demonstrates how to decode a Base64-encoded atom back into plain text using the `base64_encoded/3` predicate.
In this example, default options are used:
•	The result is returned as an atom.
•	Padding characters in the input are allowed (`padding(true)`).
•	The classic Base64 character set is used (`charset(classic)`).

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/base64.pl'),
base64_encoded(X, 'SGVsbG8gV29ybGQ=', []).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 11090
answer:
  has_more: false
  variables: ["X"]
  results:
  - substitutions:
    - variable: X
      expression: "'Hello World'"
```

### Decode a Base64 encoded atom with explicit defaults

This scenario demonstrates how to decode a Base64-encoded value back into plain text using the `base64_encoded/3` predicate.
The following options are used:

- `padding(true)` – padding characters in the input are allowed (default).
- `charset(classic)` – the classic Base64 character set is used (default).

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/base64.pl'),
base64_encoded(X, 'SGVsbG8gV29ybGQ=', [padding(true), charset(classic)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 12032
answer:
  has_more: false
  variables: ["X"]
  results:
  - substitutions:
    - variable: X
      expression: "'Hello World'"
```

### Encode text using a specific character encoding

This scenario demonstrates how the `encoding/1` option changes the bytes that are Base64-encoded before rendering the
final Base64 text.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/base64.pl'),
base64_encoded('café', X, [encoding('iso-8859-1')]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 8279
answer:
  has_more: false
  variables: ["X"]
  results:
  - substitutions:
    - variable: X
      expression: "'Y2Fm6Q=='"
```

### Error on incorrect charset option

This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid value is provided for the
`charset` option.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/base64.pl'),
base64_encoded('Hello World', X, [charset(bad)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4439
answer:
  has_more: false
  variables: ["X"]
  results:
  - error: "error(domain_error(charset,bad),base64_encoded/3)"
    substitutions:
```

### Error on incorrect padding option

This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid value is provided for the
`padding` option.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/base64.pl'),
base64_encoded('Hello World', X, [padding(bad)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4770
answer:
  has_more: false
  variables: ["X"]
  results:
  - error: "error(domain_error(padding,bad),base64_encoded/3)"
    substitutions:
```

### Error on removed as option

This scenario demonstrates how the `base64_encoded/3` predicate behaves when the removed `as` option is provided.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/base64.pl'),
base64_encoded('Hello World', X, [as(bad)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4406
answer:
  has_more: false
  variables: ["X"]
  results:
  - error: "error(type_error(option,as(bad)),base64_encoded/3)"
    substitutions:
```

### Error on incorrect encoding option

This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid value is provided for the
`encoding` option.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/base64.pl'),
base64_encoded(X, 'SGVsbG8gV29ybGQ=', [encoding(unknown)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 11778
answer:
  has_more: false
  variables: ["X"]
  results:
  - error: "error(type_error(charset,unknown),base64_encoded/3)"
    substitutions:
```

### Error on incorrect encoding option (2)

This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid type is provided for the
`encoding` option.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/base64.pl'),
base64_encoded(X, 'SGVsbG8gV29ybGQ=', [encoding(bad, 'very bad')]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4439
answer:
  has_more: false
  variables: ["X"]
  results:
  - error: "error(type_error(option,encoding(bad,very bad)),base64_encoded/3)"
    substitutions:
```

### Error on unknown option name

This scenario demonstrates how the `base64_encoded/3` predicate behaves when an unknown option name is provided.
This helps catch typos in option names (e.g., `chatset` instead of `charset`).

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/base64.pl'),
base64_encoded('Hello World', X, [chatset(classic)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4415
answer:
  has_more: false
  variables: ["X"]
  results:
  - error: "error(type_error(option,chatset(classic)),base64_encoded/3)"
    substitutions:
```
