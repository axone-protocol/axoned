---
sidebar_position: 32
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# catch/3

## Description

`catch/3` is a predicate that catches exceptions thrown during the execution of a goal.

## Signature

```text
catch(+Goal, ?`catch/3`er, +Recover)
```

Where:

- Goal is the goal to execute.
- `catch/3`er is the exception pattern to catch.
- Recover is the goal to execute when the exception is caught.
