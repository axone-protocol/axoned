---
sidebar_position: 33
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# comet_info/1

## Module

This predicate is provided by `chain.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/chain.pl').
```

## Description

Unifies CometInfo with the current Comet block info dict exposed by the VFS.

Returned term shape:

```prolog
comet{
  validators_hash: [Byte],
  proposer_address: [Byte],
  evidence: [evidence{
    type: Type,
    validator: validator{address:[Byte], power:Power},
    height: Height,
    time: Time,
    total_voting_power: TotalVotingPower
  }],
  last_commit: commit_info{
    round: Round,
    votes: [vote_info{
      block_id_flag: BlockIDFlag,
      validator: validator{address:[Byte], power:Power}
    }]
  }
}.
```

where:

- Byte is an integer in [0,255].
- Time is a Unix timestamp in seconds (0 when unset).
- Empty lists are returned when data is unavailable.

## Signature

```text
comet_info(-CometInfo) is det
```

## Examples

### Retrieve current comet block info

This scenario demonstrates how to retrieve the current Comet block information available to the query.
The Comet info contains consensus-related metadata such as proposer address, validators hash,
evidence, and last commit information.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/chain.pl'),
comet_info(CometInfo).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3991
answer:
  has_more: false
  variables: ["CometInfo"]
  results:
  - substitutions:
    - variable: CometInfo
      expression: >-
        comet{evidence:[],last_commit:commit_info{round:0,votes:[]},proposer_address:[],validators_hash:[]}
```

### Retrieve proposer address from current comet block info

This scenario demonstrates how to read a specific field from comet_info/1.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- consult('/v1/lib/chain.pl').

proposer_address(Address) :-
    comet_info(CometInfo),
    Address = CometInfo.proposer_address.
```

- **Given** the query:

```  prolog
proposer_address(Address).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3994
answer:
  has_more: false
  variables: ["Address"]
  results:
  - substitutions:
    - variable: Address
      expression: "[]"
```

### Retrieve last commit round from current comet block info

This scenario demonstrates how to read nested fields from comet_info/1.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- consult('/v1/lib/chain.pl').

last_commit_round(Round) :-
    comet_info(CometInfo),
    Round = CometInfo.last_commit.round.
```

- **Given** the query:

```  prolog
last_commit_round(Round).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3995
answer:
  has_more: false
  variables: ["Round"]
  results:
  - substitutions:
    - variable: Round
      expression: "0"
```

### Evaluate a condition based on comet evidence and commit round

This scenario demonstrates how to combine comet_info/1 fields in a rule.
It checks that there is no evidence and that the last commit round is non-negative.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- consult('/v1/lib/chain.pl').

healthy_consensus_snapshot :-
    comet_info(CometInfo),
    CometInfo.evidence = [],
    CometInfo.last_commit.round >= 0,
    !.
```

- **Given** the query:

```  prolog
healthy_consensus_snapshot.
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3991
answer:
  has_more: false
  results:
  - substitutions:
```
