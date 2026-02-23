---
sidebar_position: 152
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# setup_call_cleanup/3

## Module

Built-in predicate.

## Description

Runs Setup once, then Goal, and always executes Cleanup exactly once for
this deterministic execution path:

- on success of Goal;
- on failure of Goal;
- on exception raised by Goal (then rethrows).

This implementation is intended for deterministic goals in this runtime.

## Signature

```text
setup_call_cleanup(:Setup, :Goal, :Cleanup) is det
```
