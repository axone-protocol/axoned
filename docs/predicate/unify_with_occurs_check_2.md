---
sidebar_position: 107
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# unify_with_occurs_check/2

## Description

`unify_with_occurs_check/2` is a predicate that unifies two terms with the occurs check enabled.

## Signature

```text
unify_with_occurs_check(+Left, +Right)
```

Where:

- Left is the first term to unify.
- Right is the second term to unify.

The occurs check prevents the creation of infinite structures by ensuring that a variable does not occur in the term it is being unified with.
