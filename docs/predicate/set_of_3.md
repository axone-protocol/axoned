---
sidebar_position: 103
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# set_of/3

## Description

`set_of/3` is a predicate that collects unique solutions to a goal in sorted order.

## Signature

```text
setof(?Template, +Goal, ?Set)
```

Where:

- Template is the term to collect for each solution.
- Goal is the goal to find solutions for.
- Set is unified with the sorted list of unique solutions.
