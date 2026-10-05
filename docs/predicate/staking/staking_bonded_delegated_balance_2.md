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

## Examples

### Sum bonded delegations including a jailed bonded validator

Here is the feature setup:

- **Given** the program:

```  prolog
:- consult('/v1/lib/staking.pl').
```

- **And** the staking bond denomination is "ustake"
- **And** the staking validators:

| operator | status | jailed |
| --- | --- | --- |
| axonevaloper1bonded | bonded | false |
| axonevaloper1jailed | bonded | true |
| axonevaloper1unbonding | unbonding | false |
| axonevaloper1unbonded | unbonded | false |

Here are the steps of the scenario:

- **Given** the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following staking delegations:

| validator | denom | amount |
| --- | --- | --- |
| axonevaloper1unbonding | ustake | 900 |
| axonevaloper1jailed | ustake | 75 |
| axonevaloper1unbonded | ustake | 800 |
| axonevaloper1bonded | ustake | 25 |

- **And** the query:

```  prolog
staking_bonded_delegated_balance('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', Balance).
```

- **When** the query is run
- **Then** the logical answer we get is:

```  yaml
answer:
  has_more: false
  variables: [Balance]
  results:
  - substitutions:
    - variable: Balance
      expression: "coin(ustake,100)"
```
