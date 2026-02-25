---
sidebar_position: 7
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# assertz/1

## Description

`assertz/1` is a predicate that asserts a clause into the database as the last clause of the predicate.

## Signature

```text
assertz(+Clause)
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
assert_fact :- assertz(father(john, pete)).
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
assertz(parent(john, alice)).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3975
answer:
  has_more: false
  results:
  - error: "error(permission_error(modify,static_procedure,parent/2),assertz/1)"
```

### Show that the fact is asserted at the end of the database

This scenario demonstrates that the assertz/1 predicate adds the fact to the end of the database. This means that
the fact is the last fact to be matched when a query is run.

This is in contrast to the asserta/1 predicate, which adds the fact to the beginning of the database.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- dynamic(parent/2).

parent(jane, alice).

assert_fact :- assertz(parent(john, alice)).
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
      expression: jane
  - substitutions:
    - variable: X
      expression: john
```

### Add and remove items in an inventory

This scenario demonstrates how to maintain a dynamic list of items (like in an inventory system) by representing each item
as a fact in the Prolog knowledge base. By using dynamic predicates, we can add items to the inventory and remove them on demand.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- dynamic(inventory/1).

add_item(Item) :- assertz(inventory(Item)).
remove_item(Item) :- retract(inventory(Item)).
```

- **And** the query:

```  prolog
add_item('apple'),
add_item('banana'),
add_item('orange'),
remove_item('banana'),
findall(I, inventory(I), CurrentInventory).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3984
answer:
  has_more: false
  variables: ["I","CurrentInventory"]
  results:
  - substitutions:
    - variable: CurrentInventory
      expression: "[apple,orange]"
```
