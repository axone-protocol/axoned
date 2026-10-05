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

## Examples

### Include jailed validators if their bond status is still bonded

Here is the feature setup:

- **Given** the program:

```  prolog
:- consult('/v1/lib/staking.pl').
```

Here are the steps of the scenario:

- **Given** the staking validators:

| operator | status | jailed |
| --- | --- | --- |
| axonevaloper1bonded | bonded | false |
| axonevaloper1jailed | bonded | true |
| axonevaloper1unbonding | unbonding | false |

- **And** the query:

```  prolog
staking_bonded_validators(Validators).
```

- **When** the query is run
- **Then** the logical answer we get is:

```  yaml
answer:
  has_more: false
  variables: [Validators]
  results:
  - substitutions:
    - variable: Validators
      expression: "[validator{operator:axonevaloper1bonded,status:bonded},validator{operator:axonevaloper1jailed,status:bonded}]"
```
