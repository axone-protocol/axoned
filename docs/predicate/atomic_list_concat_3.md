---
sidebar_position: 13
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# atomic_list_concat/3

## Module

Built-in predicate.

## Description

Unifies Atom with the concatenation of the atomic textual representation of
each element in List, inserting Separator between adjacent elements.

where:

- List is a proper list of ground terms. Each element is converted using term_to_atom/2, so atoms, numbers,
  double-quoted strings, lists and compounds are supported;
- Separator is an atom inserted between adjacent elements;
- Atom is an atom representing the concatenation of the textual representation of each element in List.

Throws:

- error(instantiation_error, atomic_list_concat/3) when List or Separator is insufficiently instantiated;
- error(type_error(list, List), atomic_list_concat/3) when List is not a proper list;
- error(type_error(atom, Separator), atomic_list_concat/3) when Separator is instantiated but is not an atom.

## Signature

```text
atomic_list_concat(+List, +Separator, ?Atom) is det
```

## Examples

### Concatenate values with a separator

This scenario demonstrates how `atomic_list_concat/3` inserts a separator between the textual representation of each list element.

Here are the steps of the scenario:

- **Given** the program:

```  prolog

```

- **Given** the query:

```  prolog
atomic_list_concat([cosmos, hub, 4], '-', Atom).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4147
answer:
  has_more: false
  variables: ["Atom"]
  results:
  - substitutions:
    - variable: Atom
      expression: "'cosmos-hub-4'"
```

### Build a URI-like atom from separate parts

This scenario demonstrates how `atomic_list_concat/3` can be used to assemble a structured atom from reusable parts.

Here are the steps of the scenario:

- **Given** the program:

```  prolog

```

- **Given** the query:

```  prolog
atomic_list_concat([scheme, host, path], '://', URI).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4179
answer:
  has_more: false
  variables: ["URI"]
  results:
  - substitutions:
    - variable: URI
      expression: "'scheme://host://path'"
```
