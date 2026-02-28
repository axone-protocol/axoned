---
sidebar_position: 115
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# term_to_atom/2

## Module

Built-in predicate.

## Description

Relates a ground Term with its textual Atom representation.

where:

- Term is a ground term that unifies with the parsed representation of Atom;
- Atom is an atom containing a canonical textual representation of Term.

When Term is ground, Atom is unified with a canonical textual representation
that can be parsed back by this predicate. When Atom is instantiated, it is
parsed back into Term.

The supported syntax matches the canonical text produced here: atoms, quoted
atoms, numbers, double-quoted strings (lists of one-character atoms), lists and compounds.

Throws:

- error(instantiation_error, term_to_atom/2) when both arguments are variables;
- error(type_error(atom, Atom), term_to_atom/2) when Atom is instantiated but is not an atom;
- error(syntax_error(term), term_to_atom/2) when Atom is an atom that does not contain a valid canonical term.

## Signature

```text
term_to_atom(?Term, ?Atom) is det
```

## Examples

### Convert a ground term into a canonical atom

This scenario demonstrates how `term_to_atom/2` turns a ground term into a canonical atom that can be reused later.

Here are the steps of the scenario:

- **Given** the program:

```  prolog

```

- **Given** the query:

```  prolog
term_to_atom(greeting(hello, [world, 42]), Atom).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4172
answer:
  has_more: false
  variables: ["Atom"]
  results:
  - substitutions:
    - variable: Atom
      expression: "'greeting(hello,[world,42])'"
```

### Parse a canonical atom back into a term

This scenario demonstrates how `term_to_atom/2` reads an atom back into a Prolog term, including double-quoted strings.

Here are the steps of the scenario:

- **Given** the program:

```  prolog

```

- **Given** the query:

```  prolog
term_to_atom(Term, 'payload(\"hi\", [foo, 42])').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 5033
answer:
  has_more: false
  variables: ["Term"]
  results:
  - substitutions:
    - variable: Term
      expression: payload([h,i],[foo,42])
```
