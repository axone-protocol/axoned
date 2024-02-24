---
sidebar_position: 16
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# source_file/1

## Description

`source_file/1` is a predicate that unify the given term with the currently loaded source file.

The signature is as follows:

```text
source_file(?File).
```

Where:

- File represents a loaded source file.

## Examples

```text
# Query all the loaded source files, in alphanumeric order.
- source_file(File).

# Query the given source file is loaded.
- source_file('foo.pl').
```
