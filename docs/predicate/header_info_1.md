---
sidebar_position: 34
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# header_info/1

## Module

This predicate is provided by `chain.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/chain.pl').
```

## Description

Unifies HeaderInfo with the current SDK header info dict exposed by the VFS.

Returned term shape:

```prolog
header{
  height: Height,
  hash: [Byte],
  time: Time,
  chain_id: ChainID,
  app_hash: [Byte]
}.
```

where:

- Height is the current block height.
- Time is a Unix timestamp in seconds.
- ChainID is an atom (quoted if needed).
- Byte is an integer in [0,255].

## Signature

```text
header_info(?HeaderInfo) is det
```

## Examples

### Retrieve current SDK header info

This scenario demonstrates how to retrieve the current block header information available to the query.
The header info contains useful execution context such as the block height, block time, and chain identifier.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/chain.pl'),
header_info(HeaderInfo).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3991
answer:
  has_more: false
  variables: ["HeaderInfo"]
  results:
  - substitutions:
    - variable: HeaderInfo
      expression: >-
        header{app_hash:[],chain_id:'axone-testchain-1',hash:[],height:42,time:1712745867}
```

### Retrieve the block height of the current block

This scenario demonstrates how to read the current block height from header_info/1.

Here are the steps of the scenario:

- **Given** a block with the following header:

```  yaml
height: 100
```

- **Given** the program:

```  prolog
:- consult('/v1/lib/chain.pl').

height(Height) :-
    header_info(Header),
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
gas_used: 3994
answer:
  has_more: false
  variables: ["Height"]
  results:
  - substitutions:
    - variable: Height
      expression: "100"
```

### Retrieve the block time of the current block

This scenario demonstrates how to read the current block time from header_info/1.

Here are the steps of the scenario:

- **Given** a block with the following header:

```  yaml
time: 2024-03-04T11:03:36.000Z
```

- **Given** the program:

```  prolog
:- consult('/v1/lib/chain.pl').

time(Time) :-
    header_info(Header),
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
gas_used: 3994
answer:
  has_more: false
  variables: ["Time"]
  results:
  - substitutions:
    - variable: Time
      expression: "1709550216"
```

### Evaluate a condition based on block time and height

This scenario demonstrates how to evaluate a condition based on both block time and block height.
Specifically, it checks whether block time is greater than 1709550216 seconds
(Monday 4 March 2024 11:03:36 GMT) or block height is greater than 42.

Here are the steps of the scenario:

- **Given** a block with the following header:

```  yaml
time: 2024-03-04T11:03:37.000Z
```

- **Given** the program:

```  prolog
:- consult('/v1/lib/chain.pl').

evaluate :-
    header_info(Header),
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
gas_used: 3991
answer:
  has_more: false
  results:
  - substitutions:
```
