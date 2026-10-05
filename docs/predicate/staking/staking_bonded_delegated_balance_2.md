---
sidebar_position: 1
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# staking_bonded_delegated_balance/2

## Module

This predicate is provided by `staking.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/staking.pl').
```

## Description

Unifies Balance with coin(BondDenom, Amount), the sum of delegations to bonded validators.
Empty eligible collections return coin(BondDenom, 0), using the staking bond denomination.
Amount is an integer when it fits in int64, otherwise a decimal atom. Addition is exact,
including amounts and totals larger than int64. A mismatching coin denomination raises
domain_error(staking_bond_denom, Denom); different denominations are never added.
Unbondings, redelegation entries, rewards, and bank balances are excluded.

## Signature

```text
staking_bonded_delegated_balance(+Delegator, -Balance) is det
```
