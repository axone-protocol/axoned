---
sidebar_position: 10
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# block_header/1

## Description

`block_header/1` is a predicate which unifies the given term with the current block header.

## Signature

```text
block_header(?Header) is det
```

where:

- Header is a Dict representing the current chain header at the time of the query.

## Examples

### Retrieve the header of the current block

This scenario demonstrates how to retrieve the header of the current block and obtain some of its properties.

The header of a block carries important information about the state of the blockchain, such as basic information (chain id, the height,
time, and height), the information about the last block, hashes, and the consensus info.

The header is represented as a Prolog Dict, which is a collection of key-value pairs.

Here are the steps of the scenario:

- **Given** a block with the following header:

```  yaml
app_hash: Q0P6b2hoSUbmpCE6o6Dme4H4FBWqdcpqo89DrpBYSHQ=
chain_id: axone-localnet
height: 33
next_validators_hash: EIQFMnCDepfXD2e3OeL1QoEfmu6BZQbKR500Wkl4gK0=
proposer_address: yz7PSKMWniQlQWMd7LskBABgDKQ=
time: "2024-11-22T21:22:04.676789Z"
```

- **Given** the query:

```  prolog
block_header(Header).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 33
gas_used: 3975
answer:
  has_more: false
  variables: ["Header"]
  results:
  - substitutions:
    - variable: Header
      expression: >-
        header{app_hash:[67,67,250,111,104,104,73,70,230,164,33,58,163,160,230,123,129,248,20,21,170,117,202,106,163,207,67,174,144,88,72,116],chain_id:'axone-localnet',consensus_hash:[],data_hash:[],evidence_hash:[],height:33,last_block_id:block_id{hash:[],part_set_header:part_set_header{hash:[],total:0}},last_commit_hash:[],last_results_hash:[],next_validators_hash:[16,132,5,50,112,131,122,151,215,15,103,183,57,226,245,66,129,31,154,238,129,101,6,202,71,157,52,90,73,120,128,173],proposer_address:[203,62,207,72,163,22,158,36,37,65,99,29,236,187,36,4,0,96,12,164],time:1732310524,validators_hash:[],version:consensus{app:0,block:0}}
```

### Retrieve the block height of the current block

This scenario demonstrates how to retrieve the block height of the current block.

Here are the steps of the scenario:

- **Given** a block with the following header:

```  yaml
height: 100
```

- **Given** the program:

```  prolog
height(Height) :-
    block_header(Header),
    Height = Header.height.
```

- **Given** the query:

```  prolog
height(Height).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 100
gas_used: 3978
answer:
  has_more: false
  variables: ["Height"]
  results:
  - substitutions:
    - variable: Height
      expression: "100"
```

### Retrieve the block time of the current block

This scenario demonstrates how to retrieve the block time of the current block.

Here are the steps of the scenario:

- **Given** a block with the following header:

```  yaml
time: 2024-03-04T11:03:36.000Z
```

- **Given** the program:

```  prolog
time(Time) :-
    block_header(Header),
    Time = Header.time.
```

- **Given** the query:

```  prolog
time(Time).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3978
answer:
  has_more: false
  variables: ["Time"]
  results:
  - substitutions:
    - variable: Time
      expression: "1709550216"
```

### Evaluate a condition based on block time and height

This scenario demonstrates how to evaluate a condition that depends on both the block time and block height.
Specifically, it checks whether the block time is greater than 1709550216 seconds (Monday 4 March 2024 11:03:36 GMT)
or the block height is greater than 42.

Here are the steps of the scenario:

- **Given** a block with the following header:

```  yaml
time: 2024-03-04T11:03:37.000Z
```

- **Given** the program:

```  prolog
evaluate :-
    block_header(Header),
    (Header.time > 1709550216; Header.height > 42),
    !.
```

- **Given** the query:

```  prolog
evaluate.
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3981
answer:
  has_more: false
  results:
  - substitutions:
```
