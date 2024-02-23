---
sidebar_position: 14
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# open/4

## Description

`open/4` is a predicate that unify a stream with a source sink on a virtual file system.

The signature is as follows:

```text
open(+SourceSink, +Mode, ?Stream, +Options)
```

Where:

- SourceSink is an atom representing the source or sink of the stream. The atom typically represents a resource that can be opened, such as a URI. The URI scheme determines the type of resource that is opened.
- Mode is an atom representing the mode of the stream \(read, write, append\).
- Stream is the stream to be opened.
- Options is a list of options. No options are currently defined, so the list should be empty.

## Examples

```text
# `open/4` a stream from a cosmwasm query.
# The Stream should be read as a string with a read_string/3 predicate, and then closed with the close/1 predicate.
- open('cosmwasm:okp4-objectarium:okp412kgx?query=%7B%22object_data%22%3A%7B%...4dd539e3%22%7D%7D', 'read', Stream, [])
```
