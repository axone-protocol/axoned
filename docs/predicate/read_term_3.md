---
sidebar_position: 83
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# read_term/3

## Description

`read_term/3` is a predicate that reads a term from a stream or alias.

The signature is as follows:

```text
read_term(+Stream, -Term, +Options)
```

where:

- Stream represents the stream or alias to read the term from.
- Term represents the term to read.
- Options represents the options to control the reading process.

Valid options are:

- singletons\(Vars\): Vars is unified with a list of variables that occur only once in the term.
- variables\(Vars\): Vars is unified with a list of variables that occur in the term.
- variable\_names\(Vars\): Vars is unified with a list of Name = Var terms, where Name is an atom and Var is a variable.
