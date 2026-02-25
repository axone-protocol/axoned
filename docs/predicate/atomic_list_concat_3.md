---
sidebar_position: 13
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# atomic_list_concat/3

## Description

`atomic_list_concat/3` is a predicate that unifies an Atom with the concatenated elements of a List using a given separator.

The atomic\_list\_concat/3 predicate creates an atom just like atomic\_list\_concat/2, but inserts Separator between each pair of inputs.

## Signature

```text
atomic_list_concat(+List, +Separator, ?Atom)
```

where:

- List is a list of strings, atoms, integers, floating point numbers or non\-integer rationals
- Separator is an atom \(possibly empty\)
- Atom is an Atom representing the concatenation of the elements of List
