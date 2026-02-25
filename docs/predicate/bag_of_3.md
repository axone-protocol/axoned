---
sidebar_position: 14
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# bag_of/3

## Description

`bag_of/3` is a predicate that collects solutions to a goal grouped by free variables.

## Signature

```text
bagof(?Template, +Goal, ?Bag)
```

Where:

- Template is the term to collect for each solution.
- Goal is the goal to find solutions for.
- Bag is unified with the list of collected solutions.

Unlike findall/3, bagof/3 fails if there are no solutions and groups solutions by free variables.
