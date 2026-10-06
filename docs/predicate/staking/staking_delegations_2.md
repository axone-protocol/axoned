---
sidebar_position: 3
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# staking_delegations/2

## Module

This predicate is provided by `staking.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/staking.pl').
```

## Description

Unifies Delegations with all staking delegations of Delegator, regardless of validator status.
Delegator must be an AXONE account address atom in Bech32 format.

Each item has the shape:

```prolog
delegation{validator: Validator, balance: coin(Denom, Amount)}
```

Validator and Denom are atoms. Amount is an integer when it fits in int64,
otherwise an atom preserving its full decimal precision.

## Signature

```text
staking_delegations(+Delegator, -Delegations) is det
```

## Examples

### Query all delegations of an account

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- consult('/v1/lib/staking.pl').
```

- **And** the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following staking delegations:

| validator | denom | amount |
| --- | --- | --- |
| axonevaloper1validator | uaxone | 1250000 |

- **And** the staking validators:

| operator | status | jailed |
| --- | --- | --- |
| axonevaloper1validator | unbonded | false |

- **Given** the query:

```  prolog
staking_delegations('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', Delegations).
```

- **When** the query is run
- **Then** the logical answer we get is:

```  yaml
answer:
  has_more: false
  variables: ["Delegations"]
  results:
  - substitutions:
    - variable: Delegations
      expression: "[delegation{balance:coin(uaxone,1250000),validator:axonevaloper1validator}]"
```
