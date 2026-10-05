---
sidebar_position: 7
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# staking_validators/2

## Module

This predicate is provided by `staking.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/staking.pl').
```

## Description

Unifies Validators with validators having bond Status: bonded, unbonding, or unbonded.
Status must be an atom. Each item is validator{operator: Operator, status: Status}.
Operator is a validator operator address atom. Bond status is independent of jailed status:
bonded means membership in the active staking validator set, not an unjailed guarantee.

## Signature

```text
staking_validators(+Status, -Validators) is det
```
