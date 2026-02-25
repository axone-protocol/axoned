---
sidebar_position: 127
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# write_term/3

## Description

`write_term/3` is a predicate that writes a term to a stream or alias.

The signature is as follows:

```text
write_term(+Stream, +Term, +Options)
```

where:

- Stream represents the stream or alias to write the term to.
- Term represents the term to write.
- Options represents the options to control the writing process.

Valid options are:

- quoted\(Bool\): If true, atoms and strings that need quotes will be quoted. The default is false.
- ignore\_ops\(Bool\): If true, the generic term representation \(\<functor\>\(\<args\> ... \)\) will be used for all terms. Otherwise \(default\), operators will be used where appropriate.
- numbervars\(Bool\): If true, variables will be numbered. The default is false.
- variable\_names\(\+List\): Assign names to variables in Term. List is a list of Name = Var terms, where Name is an atom and Var is a variable.
- max\_depth\(\+Int\): The maximum depth to which the term is written. The default is infinite.
