---
sidebar_position: 12
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# atomic_list_concat/2

## Module

Built-in predicate.

## Description

Unifies Atom with the concatenation of the atomic textual representation of
each element in List.

where:

- List is a proper list of ground terms. Each element is converted using term_to_atom/2, so atoms, numbers,
  double-quoted strings, lists and compounds are supported;
- Atom is an atom representing the concatenation of the textual representation of each element in List.

Throws:

- error(instantiation_error, atomic_list_concat/2) when List is insufficiently instantiated;
- error(type_error(list, List), atomic_list_concat/2) when List is not a proper list.

## Signature

```text
atomic_list_concat(+List, ?Atom) is det
```

## Examples

### Concatenate atomic values into a single atom

This scenario demonstrates how `atomic_list_concat/2` concatenates the textual representation of several atomic values.

Here are the steps of the scenario:

- **Given** the program:

```  prolog

```

- **Given** the query:

```  prolog
atomic_list_concat([hello, '-', 42, '-', world], Atom).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4295
answer:
  has_more: false
  variables: ["Atom"]
  results:
  - substitutions:
    - variable: Atom
      expression: "'hello-42-world'"
```
