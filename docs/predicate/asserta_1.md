---
sidebar_position: 2
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# asserta/1

## Description

`asserta/1` is a predicate that asserts a clause into the database as the first clause of the predicate.

## Signature

```text
asserta(+Clause)
```

Where:

- Clause is the clause to assert into the database.

## Examples

### Assert a fact into the database

This scenario demonstrates the process of asserting a new fact into a Prolog database. In Prolog, asserting a fact means
adding a new piece of information or *knowledge* into the database, allowing it to be referenced in subsequent queries.
This is particularly useful when you want to dynamically extend the knowledge base with facts or rules based on conditions
or interactions during runtime.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
assert_fact :- asserta(father(john, pete)).
```

- **Given** the query:

```  prolog
assert_fact, father(X, Y).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3977
answer:
  has_more: false
  variables: ["X", "Y"]
  results:
  - substitutions:
    - variable: X
      expression: john
    - variable: 'Y'
      expression: pete
```

### Only dynamic predicates can be asserted

This scenario demonstrates that only dynamic predicates can be asserted. In Prolog, dynamic predicates are those that can be
modified during runtime. This is in contrast to static predicates, which are fixed and cannot be modified.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
parent(jane, alice).
```

- **Given** the query:

```  prolog
asserta(parent(john, alice)).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3975
answer:
  has_more: false
  results:
  - error: "error(permission_error(modify,static_procedure,parent/2),asserta/1)"
```

### Show that the fact is asserted at the beginning of the database

This scenario demonstrates that the asserta/1 predicate adds the fact to the beginning of the database. This means that
the fact is the first fact to be matched when a query is run.

This is in contrast to the assertz/1 predicate, which adds the fact to the end of the database.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- dynamic(parent/2).

parent(jane, alice).

assert_fact :- asserta(parent(john, alice)).
```

- **Given** the query:

```  prolog
assert_fact, parent(X, alice).
```

- **When** the query is run (limited to 2 solutions)
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3977
answer:
  has_more: false
  variables: ["X"]
  results:
  - substitutions:
    - variable: X
      expression: john
  - substitutions:
    - variable: X
      expression: jane
```

### Shows a simple counter example

This scenario demonstrates a simple counter example using the `asserta/1` and `retract/1` predicates.
In this example, we represent the value of the counter as a dynamic predicate `counter/1` that is asserted and retracted
to each time the value of the counter is incremented or decremented.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- dynamic(counter/1).

counter(0).

increment_counter :- retract(counter(X)), Y is X + 1, asserta(counter(Y)).
decrement_counter :- retract(counter(X)), Y is X - 1, asserta(counter(Y)).
```

- **Given** the query:

```  prolog
counter(InitialValue), increment_counter, increment_counter, counter(IncrementedValue), decrement_counter, counter(DecrementedValue).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3989
answer:
  has_more: false
  variables: ["InitialValue", "IncrementedValue", "DecrementedValue"]
  results:
  - substitutions:
    - variable: InitialValue
      expression: 0
    - variable: IncrementedValue
      expression: 2
    - variable: DecrementedValue
      expression: 1
```
