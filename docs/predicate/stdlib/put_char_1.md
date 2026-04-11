---
sidebar_position: 71
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# put_char/1

## Module

Built-in predicate.

## Description

Writes Char to the current output stream.

## Signature

```text
put_char(+Char) is det
```

## Examples

### Write a single character to user output

This scenario demonstrates using put_char/1 to write a single character to the current output stream.
The character appears in the user_output field of the response.

Here are the steps of the scenario:

- **Given** the module configuration:

```  json
{
  "limits": {
    "max_user_output_size": 10
  }
}
```

- **Given** the query:

```  prolog
put_char('b').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3986
answer:
  has_more: false
  variables:
  results:
  - substitutions:
user_output: "b"
```

### Write multiple characters to user output

This scenario demonstrates chaining multiple put_char/1 calls to write several characters.
Each character is appended to the user output stream.

Here are the steps of the scenario:

- **Given** the module configuration:

```  json
{
  "limits": {
    "max_user_output_size": 10
  }
}
```

- **Given** the query:

```  prolog
put_char('a'), put_char('b'), put_char('c').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4042
answer:
  has_more: false
  variables:
  results:
  - substitutions:
user_output: "abc"
```

### Write characters with user output size limit

This scenario shows how the user output is truncated when it exceeds the configured max_user_output_size limit.
The limit is measured in bytes, so only the last bytes that fit within the limit are kept.

Here are the steps of the scenario:

- **Given** the module configuration:

```  json
{
  "limits": {
    "max_user_output_size": 3
  }
}
```

- **Given** the query:

```  prolog
put_char('h'), put_char('e'), put_char('l'), put_char('l'), put_char('o').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4098
answer:
  has_more: false
  variables:
  results:
  - substitutions:
user_output: "llo"
```

### Write UTF-8 character

This scenario illustrates writing UTF-8 characters using put_char/1.
Multi-byte characters like emojis occupy more space in the buffer.

Here are the steps of the scenario:

- **Given** the module configuration:

```  json
{
  "limits": {
    "max_user_output_size": 10
  }
}
```

- **Given** the program:

```  prolog
log_message([]).
log_message([H|T]) :-
    put_char(H),
    log_message(T).
```

- **Given** the query:

```  prolog
log_message("😀").
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4093
answer:
  has_more: false
  variables:
  results:
  - substitutions:
user_output: "😀"
```
