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

## Examples

### Keep balances and order for bonded validators only

Here is the feature setup:

- **Given** the program:

```  prolog
:- consult('/v1/lib/staking.pl').
```

- **And** the staking validators:

| operator | status | jailed |
| --- | --- | --- |
| axonevaloper1bonded | bonded | false |
| axonevaloper1jailed | bonded | true |
| axonevaloper1unused | bonded | false |
| axonevaloper1unbonding | unbonding | false |
| axonevaloper1unbonded | unbonded | false |

Here are the steps of the scenario:

- **Given** the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following staking delegations:

| validator | denom | amount |
| --- | --- | --- |
| axonevaloper1unbonding | uaxone | 900 |
| axonevaloper1jailed | uaxone | 75 |
| axonevaloper1unbonded | uaxone | 800 |
| axonevaloper1bonded | uaxone | 25 |
| axonevaloper1absent | uaxone | 700 |

- **And** the query:

```  prolog
staking_delegations_to_bonded_validators('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', Delegations).
```

- **When** the query is run
- **Then** the logical answer we get is:

```  yaml
answer:
  has_more: false
  variables: [Delegations]
  results:
  - substitutions:
    - variable: Delegations
      expression: "[delegation{balance:coin(uaxone,75),validator:axonevaloper1jailed},delegation{balance:coin(uaxone,25),validator:axonevaloper1bonded}]"
```
