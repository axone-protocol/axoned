---
sidebar_position: 12
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# current_output/1

## Description

`current_output/1` is a predicate that unifies the given term with the current output stream.

## Signature

```text
current_output(-Stream) is det
```

where:

- Stream represents the current output stream.

This predicate connects to the default output stream available for user interactions, allowing the user to perform write operations.

The outcome of the stream's content throughout the execution of a query is provided as a string within the user\_output field in the query's response. However, it's important to note that the maximum length of the output is constrained by the max\_query\_output\_size setting, meaning only the final max\_query\_output\_size bytes \(not characters\) of the output are included in the response.

## Examples

### Write a char to the current output

This scenario demonstrates how to write a character to the current output, and get the content in the response of the
request.

Here are the steps of the scenario:

- **Given** the module configuration:

```  json
{
  "limits": {
    "max_user_output_size": 5
  }
}
```

- **Given** the program:

```  prolog
write_char_to_user_output(C) :-
    current_output(UserStream), % get the current output stream
    put_char(UserStream, C).    % write the char to the user stream
```

- **Given** the query:

```  prolog
write_char_to_user_output(x).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4043
answer:
  has_more: false
  variables:
  results:
  - substitutions:
user_output: |
  x
```

### Write characters to the current output (without limit)

This scenario demonstrates how to write some characters to the current output, and get the content in the response of the
request. This is helpful for debugging purposes.

Here are the steps of the scenario:

- **Given** the module configuration:

```  json
{
  "limits": {
    "max_user_output_size": 15
  }
}
```

- **Given** the program:

```  prolog
log_message(Message) :-
    current_output(UserStream), % get the current output stream
    write(UserStream, Message), % write the message to the user stream
    put_char(UserStream, '\n').
```

- **Given** the query:

```  prolog
log_message('Hello world!').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4045
answer:
  has_more: false
  variables:
  results:
  - substitutions:
user_output: |
  Hello world!

```

### Write characters to the current output (with limit)

This scenario demonstrates the process of writing characters to the current user output, with a limit configured
in the logic module. So if the message is longer than this limit, the output will be truncated.

Here are the steps of the scenario:

- **Given** the module configuration:

```  json
{
  "limits": {
    "max_user_output_size": 5
  }
}
```

- **Given** the program:

```  prolog
log_message(Message) :-
    current_output(UserStream), % get the current output stream
    write(UserStream, Message). % write the message to the user stream
```

- **Given** the query:

```  prolog
log_message('Hello world!').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4044
answer:
  has_more: false
  variables:
  results:
  - substitutions:
user_output: |
  orld!
```

### Write UTF-8 character to the current output (with limit)

This scenario illustrates the impact of UTF-8 characters on output limits measured in bytes, not character count.
Characters such as emojis require more space; for example, the wizard emoji (🧙) occupies 4 bytes, effectively counting
as four units. As a result, the limit is reached more quickly with these characters, which means that the number of
characters in the user output is less than expected.

Here are the steps of the scenario:

- **Given** the module configuration:

```  json
{
  "limits": {
    "max_user_output_size": 5
  }
}
```

- **Given** the program:

```  prolog
log_message([]).
log_message([H|T]) :-
    current_output(UserStream),
    put_char(UserStream, H),
    log_message(T).
```

- **Given** the query:

```  prolog
log_message("Hello 🧙!").
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4065
answer:
  has_more: false
  variables:
  results:
  - substitutions:
user_output: "🧙!"
```
