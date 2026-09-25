---
sidebar_position: 1
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

Unifies Delegations with the active staking delegations of Delegator.
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

### Query active delegations of an account

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- consult('/v1/lib/staking.pl').
```

- **And** the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following staking delegations:

| key | value |
| --- | ----- |
| validator | denom |
| axonevaloper1validator | uaxone |

- **Given** the query:

```  prolog
staking_delegations('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', Delegations).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 9022
answer:
  has_more: false
  variables: ["Delegations"]
  results:
  - substitutions:
    - variable: Delegations
      expression: "[delegation{balance:coin(uaxone,1250000),validator:axonevaloper1validator}]"
```
