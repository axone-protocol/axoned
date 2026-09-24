---
sidebar_position: 3
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# staking_redelegations/2

## Module

This predicate is provided by `staking.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/staking.pl').
```

## Description

Unifies Redelegations with the in-flight redelegations of Delegator.
Delegator must be an AXONE account address atom in Bech32 format.

Each item has the shape:

```prolog
redelegation{
  source_validator: SourceValidator,
  destination_validator: DestinationValidator,
  entries: [redelegation_entry{balance: coin(Denom, Amount), completion_time: Time}, ...]
}
```

Time is a Unix timestamp in seconds.

## Signature

```text
staking_redelegations(+Delegator, -Redelegations) is det
```
