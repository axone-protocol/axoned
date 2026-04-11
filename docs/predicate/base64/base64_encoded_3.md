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

Relates a text value to its Base64-encoded representation as specified by
[RFC 4648](https://rfc-editor.org/rfc/rfc4648.html).

The predicate follows a functional direction:

- when `Plain` is instantiated, it encodes `Plain` into `Encoded`;
- otherwise, when `Encoded` is instantiated, it decodes `Encoded` into `Plain`;
- otherwise, it throws `instantiation_error`.

`Plain` may be an atom, a list of characters, or a list of character codes.
`Encoded` may be an atom, a list of characters, or a list of character codes.

Supported options are:

- `charset(+Charset)` where `Charset` is `classic` (default) or `url`;
- `padding(+Boolean)` where `Boolean` is `true` (default) or `false`;
- `as(+Type)` where `Type` is `string` (default) or `atom`;
- `encoding(+Encoding)` to translate between text and bytes, defaulting to `utf8`.

## Signature

```text
base64_encoded(+Plain, -Encoded, +Options) is det
```

## Examples

### Encode a string into a Base64 encoded string (with default options)

This scenario demonstrates how to encode a plain string into its Base64 representation using the `base64_encoded/3`
predicate. The default options are used, meaning:

- The output is returned as a list of characters (`as(string)`).
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
gas_used: 9043
answer:
  has_more: false
  variables: ["X"]
  results:
  - substitutions:
    - variable: X
      expression: "['S','G','V',s,b,'G','8',g,'V','2','9',y,b,'G','Q',=]"
```

### Encode a string into a Base64 encoded atom

This scenario demonstrates how to encode a plain string into a Base64-encoded atom using the `base64_encoded/3`
predicate. The `as(atom)` option is specified, so the result is returned as a Prolog atom instead of a character
list. All other options use their default values.

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
gas_used: 9408
answer:
  has_more: false
  variables: ["X"]
  results:
  - substitutions:
    - variable: X
      expression: "'SGVsbG8gV29ybGQ='"
```

### Encode a string into a Base64 encoded atom without padding

This scenario demonstrates how to encode a plain string into a Base64-encoded atom using the `base64_encoded/3` predicate
with custom options. The following options are used:

- `as(atom)` – the result is returned as a Prolog atom.
- `padding(false)` – padding characters (`=`) are omitted.
- The classic Base64 character set is used by default (`charset(classic)`).

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/base64.pl'),
base64_encoded('Hello World', X, [as(atom), padding(false)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 10063
answer:
  has_more: false
  variables: ["X"]
  results:
  - substitutions:
    - variable: X
      expression: "'SGVsbG8gV29ybGQ'"
```

### Encode a String into a Base64 encoded atom in URL-Safe mode

This scenario demonstrates how to encode a plain string into a Base64-encoded atom using the `base64_encoded/3` predicate
with URL-safe encoding. The following options are used:

- `as(atom)` – the result is returned as a Prolog atom.
- `charset(url)` – the URL-safe Base64 alphabet is used (e.g., `-` and `_` instead of `+` and `/`).
- Padding characters are included by default (`padding(true)`).

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/base64.pl'),
base64_encoded('<<???>>', Classic, [as(atom), charset(classic)]),
base64_encoded('<<???>>', UrlSafe, [as(atom), charset(url)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 13901
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

### Decode a Base64 encoded String into plain text

This scenario demonstrates how to decode a Base64-encoded value back into plain text using the `base64_encoded/3` predicate.
The encoded input can be provided as a character list or an atom. In this example, default options are used:
•	The result (plain text) is returned as a character list (`as(string)`).
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
gas_used: 10809
answer:
  has_more: false
  variables: ["X"]
  results:
  - substitutions:
    - variable: X
      expression: "['H',e,l,l,o,' ','W',o,r,l,d]"
```

### Decode a Base64 Encoded string into a plain atom

This scenario demonstrates how to decode a Base64-encoded value back into plain text using the `base64_encoded/3` predicate,
with the result returned as a Prolog atom. The following options are used:

- `as(atom)` – the decoded plain text is returned as an atom.
- `padding(true)` – padding characters in the input are allowed (default).
- `charset(classic)` – the classic Base64 character set is used (default).

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/base64.pl'),
base64_encoded(X, 'SGVsbG8gV29ybGQ=', [as(atom)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 11678
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
base64_encoded('café', X, [as(atom), encoding('iso-8859-1')]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 8016
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
gas_used: 4445
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
gas_used: 4801
answer:
  has_more: false
  variables: ["X"]
  results:
  - error: "error(domain_error(padding,bad),base64_encoded/3)"
    substitutions:
```

### Error on incorrect as option

This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid value is provided for the
`as` option.

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
gas_used: 5253
answer:
  has_more: false
  variables: ["X"]
  results:
  - error: "error(domain_error(as,bad),base64_encoded/3)"
    substitutions:
```

### Error on incorrect encoding option

This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid value is provided for the
`encoding` option.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/base64.pl'),
base64_encoded(X, 'SGVsbG8gV29ybGQ=', [as(atom), encoding(unknown)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 12449
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
gas_used: 4475
answer:
  has_more: false
  variables: ["X"]
  results:
  - error: "error(type_error(option,encoding(bad,very bad)),base64_encoded/3)"
    substitutions:
```
