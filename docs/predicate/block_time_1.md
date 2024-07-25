---
sidebar_position: 6
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# block_time/1

## Description

`block_time/1` is a predicate which unifies the given term with the current block time.

## Signature

```text
block_time(?Time) is det
```

where:

- Time represents the current chain time at the time of the query.

## Examples

### Retrieve the block time of the current block

This scenario demonstrates how to retrieve the block time of the current block.

Here are the steps of the scenario:

- **Given** a block with the following header:

| key | value |
| --- | ----- |
| Time | 1709550216 |

- **Given** the query:

```  prolog
block_time(Time).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4140
answer:
  has_more: false
  variables: ["Time"]
  results:
  - substitutions:
    - variable: Time
      expression: "1709550216"
```

### Check that the block time is greater than a certain time

This scenario demonstrates how to check that the block time is greater than 1709550216 seconds (Monday 4 March 2024 11:03:36 GMT)
using the `block_time/1` predicate. This predicate is useful for governance which requires a certain block time to be
reached before a certain action is taken.

Here are the steps of the scenario:

- **Given** a block with the following header:

| key | value |
| --- | ----- |
| Time | 1709550217 |

- **Given** the query:

```  prolog
block_time(Time),
Time > 1709550216.
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4141
answer:
  has_more: false
  variables: ["Time"]
  results:
  - substitutions:
    - variable: Time
      expression: "1709550217"
```
