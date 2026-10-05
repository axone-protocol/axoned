---
sidebar_position: 4
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# staking_delegations_to_bonded_validators/2

## Module

This predicate is provided by `staking.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/staking.pl').
```

## Description

Unifies Delegations with the delegations of Delegator whose validators have bond status
bonded. It preserves delegation order and balances. It does not additionally filter jailed
validators. Unbondings, redelegation entries, and distribution rewards are not included.
Delegator must be an AXONE account address atom in Bech32 format.

## Signature

```text
staking_delegations_to_bonded_validators(+Delegator, -Delegations) is det
```
