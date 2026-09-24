---
sidebar_position: 4
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# staking_unbonding_delegations/2

## Module

This predicate is provided by `staking.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/staking.pl').
```

## Description

Unifies Unbondings with the unbonding staking delegations of Delegator.
Delegator must be an AXONE account address atom in Bech32 format.

Each item has the shape:

```prolog
unbonding_delegation{
  validator: Validator,
  entries: [unbonding_entry{balance: coin(Denom, Amount), completion_time: Time}, ...]
}
```

Time is a Unix timestamp in seconds. Balance is the amount currently expected
at completion, after any applicable slashing.

## Signature

```text
staking_unbonding_delegations(+Delegator, -Unbondings) is det
```
