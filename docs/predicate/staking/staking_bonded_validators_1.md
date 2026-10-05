---
sidebar_position: 2
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# staking_bonded_validators/1

## Module

This predicate is provided by `staking.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/staking.pl').
```

## Description

Unifies Validators with validators in BOND_STATUS_BONDED, including any jailed validators
whose bond status is still bonded. This is staking_validators(bonded, Validators).

## Signature

```text
staking_bonded_validators(-Validators) is det
```
