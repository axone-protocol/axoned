---
sidebar_position: 9
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# base64_encoded/3

## Description

`base64_encoded/3` is a predicate that unifies a string to a base64 encoded string as specified by [RFC 4648](<https://rfc-editor.org/rfc/rfc4648.html>).

The signature is as follows:

```text
base64_encoded(+Plain, -Encoded, +Options) is det
base64_encoded(-Plain, +Encoded, +Options) is det
base64_encoded(+Plain, +Encoded, +Options) is det
```

Where:

- Plain is an atom, a list of character codes, or list of characters containing the unencoded \(plain\) text.
- Encoded is an atom or string containing the base64 encoded text.
- Options is a list of options that can be used to control the encoding process.

## Options

The following options are supported:

- padding\(\+Boolean\)

If true \(default\), the output is padded with = characters.

- charset\(\+Charset\)

Define the encoding character set to use. The \(default\) 'classic' uses the classical rfc2045 characters. The value 'url' uses URL and file name friendly characters.

- as\(\+Type\)

Defines the type of the output. One of string \(default\) or atom.

- encoding\(\+Encoding\)

Encoding to use for translation between \(Unicode\) text and bytes \(Base64 is an encoding for bytes\). Default is utf8.

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
base64_encoded('Hello World', X, []).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3975
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
base64_encoded('Hello World', X, [as(atom)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3975
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
base64_encoded('Hello World', X, [as(atom), padding(false)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3975
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
base64_encoded('<<???>>', Classic, [as(atom), charset(classic)]),
base64_encoded('<<???>>', UrlSafe, [as(atom), charset(url)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3976
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
base64_encoded(X, 'SGVsbG8gV29ybGQ=', []).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3975
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
base64_encoded(X, 'SGVsbG8gV29ybGQ=', [as(atom)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3975
answer:
  has_more: false
  variables: ["X"]
  results:
  - substitutions:
    - variable: X
      expression: "'Hello World'"
```

### Error on incorrect charset option

This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid value is provided for the
`charset` option.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
base64_encoded('Hello World', X, [charset(bad)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3975
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
base64_encoded('Hello World', X, [padding(bad)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3975
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
base64_encoded('Hello World', X, [as(bad)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3975
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
base64_encoded(X, 'SGVsbG8gV29ybGQ=', [as(atom), encoding(unknown)]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3975
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
base64_encoded(X, 'SGVsbG8gV29ybGQ=', [encoding(bad, 'very bad')]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3975
answer:
  has_more: false
  variables: ["X"]
  results:
  - error: "error(type_error(option,encoding(bad,very bad)),base64_encoded/3)"
    substitutions:
```
