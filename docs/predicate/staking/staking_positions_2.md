---
sidebar_position: 2
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# staking_positions/2

## Module

This predicate is provided by `staking.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/staking.pl').
```

## Description

Unifies Positions with all staking-module positions of Delegator.
Delegator must be an AXONE account address atom in Bech32 format.

Positions has the shape:

```prolog
staking{
  delegations: [delegation{...}, ...],
  unbonding_delegations: [unbonding_delegation{...}, ...],
  redelegations: [redelegation{...}, ...]
}
```

This predicate covers staking-module state only. Delegation rewards belong to
the distribution module and are not included.

## Signature

```text
staking_positions(+Delegator, -Positions) is det
```
