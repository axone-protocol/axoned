---
sidebar_position: 39
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# consult/1

## Description

`consult/1` is a predicate which read files as Prolog source code.

## Signature

```text
consult(+Files) is det
```

where:

- Files represents the source files to be loaded. It can be an atom or a list of atoms representing the source files.

The Files argument are typically URIs that point to the sources file to be loaded through the Virtual File System \(VFS\). Please refer to the open/4 predicate for more information about the VFS.

## Examples

### Consult a Prolog program from the embedded library

This scenario demonstrates how to load a library file and use one of its predicates.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/lists.pl'),
member(Who, [alice,bob]).
```

- **When** the query is run (limited to 1 solutions)
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4013
answer:
  has_more: true
  variables: ["Who"]
  results:
  - substitutions:
    - variable: Who
      expression: "alice"
```

### Consult several Prolog programs at once

This scenario demonstrates consult/1 with a list of files.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- consult([
  '/v1/lib/bank.pl',
  '/v1/lib/chain.pl'
]).
```

- **Given** the query:

```  prolog
current_predicate(bank_balances/2),
current_predicate(header_info/1).
```

- **When** the query is run (limited to 2 solutions)
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4792
answer:
  has_more: false
  variables:
  results:
  - substitutions:
```
