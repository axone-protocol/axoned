---
sidebar_position: 103
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# write/1

## Module

Built-in predicate.

## Description

Writes Term to the current output stream.

## Signature

```text
write(+Term) is det
```

## Examples

### Write a simple atom to user output

This scenario demonstrates using write/1 to write an atom to the current output stream.
The term is written in human-readable form to the user_output field of the response.

Here are the steps of the scenario:

- **Given** the module configuration:

```  json
{
  "limits": {
    "max_user_output_size": 20
  }
}
```

- **Given** the query:

```  prolog
write(hello).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4002
answer:
  has_more: false
  variables:
  results:
  - substitutions:
user_output: "hello"
```

### Write a string to user output

This scenario demonstrates writing an atom containing spaces using write/1.

Here are the steps of the scenario:

- **Given** the module configuration:

```  json
{
  "limits": {
    "max_user_output_size": 30
  }
}
```

- **Given** the query:

```  prolog
write('hello world').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4010
answer:
  has_more: false
  variables:
  results:
  - substitutions:
user_output: "hello world"
```

### Write multiple terms to user output

This scenario demonstrates chaining multiple write/1 calls to output several terms.
Each term is appended directly to the user output stream without spaces.

Here are the steps of the scenario:

- **Given** the module configuration:

```  json
{
  "limits": {
    "max_user_output_size": 30
  }
}
```

- **Given** the query:

```  prolog
write('hello'), write(' '), write('world'), write('!').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4134
answer:
  has_more: false
  variables:
  results:
  - substitutions:
user_output: "hello world!"
```

### Write with user output size limit

This scenario shows how write/1 respects the max_user_output_size limit.
If the output exceeds the limit, only the last bytes that fit are retained.

Here are the steps of the scenario:

- **Given** the module configuration:

```  json
{
  "limits": {
    "max_user_output_size": 5
  }
}
```

- **Given** the query:

```  prolog
write('hello world'), put_char('!').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4038
answer:
  has_more: false
  variables:
  results:
  - substitutions:
user_output: "orld!"
```

### Write numbers and complex terms

This scenario demonstrates that write/1 can output numbers and complex terms.

Here are the steps of the scenario:

- **Given** the module configuration:

```  json
{
  "limits": {
    "max_user_output_size": 30
  }
}
```

- **Given** the query:

```  prolog
write(42), write(' '), write([1,2,3]).
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
user_output: "42 [1,2,3]"
```

### Combine write/1 and put_char/1 for formatted output

This scenario shows how write/1 and put_char/1 work together to create formatted output,
useful for debugging and logging in Prolog programs.

Here are the steps of the scenario:

- **Given** the module configuration:

```  json
{
  "limits": {
    "max_user_output_size": 51
  }
}
```

- **Given** the program:

```  prolog
log(Message) :- write('LOG: '), write(Message).
```

- **Given** the query:

```  prolog
log('Starting query'), put_char('\n'),
log('Processing data'), put_char('\n'),
log('Done').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4324
answer:
  has_more: false
  variables:
  results:
  - substitutions:
user_output: |
  LOG: Starting query
  LOG: Processing data
  LOG: Done
```
