---
sidebar_position: 29
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# read_string/3

## Description

`read_string/3` is a predicate that reads characters from the provided Stream and unifies them with String. Users can optionally specify a maximum length for reading; if the stream reaches this length, the reading stops. If Length remains unbound, the entire Stream is read, and upon completion, Length is unified with the count of characters read.

The signature is as follows:

```text
read_string(+Stream, ?Length, -String) is det
```

Where:

- Stream is the input stream to read from.
- Length is the optional maximum number of characters to read from the Stream. If unbound, denotes the full length of Stream.
- String is the resultant string after reading from the Stream.

## Examples

```text
# Given a file `foo.txt` that contains `Hello World`:

 file_to_string(File, String, Length) :-
 open(File, read, In),
 read_string(In, Length, String),
 close(Stream).

# It gives:
?- file_to_string('path/file/foo.txt', String, Length).

String = 'Hello World'
Length = 11
```
