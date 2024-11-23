---
sidebar_position: 31
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# term_to_atom/2

## Description

`term_to_atom/2` is a predicate that describes Atom as a term that unifies with Term.

## Signature

```text
term_to_atom(?Term, ?Atom)
```

where:

- Term is a term that unifies with Atom.
- Atom is an atom.

When Atom is instantiated, Atom is parsed and the result unified with Term. If Atom has no valid syntax, a syntax\_error exception is raised. Otherwise, Term is “written” on Atom using write\_term/2 with the option quoted\(true\).

## Example

```text
# Convert the atom to a term.
- term_to_atom(foo, foo).
```
