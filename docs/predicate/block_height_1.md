---
sidebar_position: 7
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# block_height/1

## Description

`block_height/1` is a predicate which unifies the given term with the current block height.

## Signature

```text
block_height(?Height) is det
```

where:

- Height represents the current chain height at the time of the query.

## Examples

### Retrieve the block height of the current block

This scenario demonstrates how to retrieve the block height of the current block.

Here are the steps of the scenario:

- **Given** a block with the following header:

| key | value |
| --- | ----- |
| Height | 100 |

- **Given** the query:

```  prolog
block_height(Height).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 100
gas_used: 3612
answer:
  has_more: false
  variables: ["Height"]
  results:
  - substitutions:
    - variable: Height
      expression: "100"
```

### Check that the block height is greater than a certain value

This scenario demonstrates how to check that the block height is greater than 100. This predicate is useful for
governance which requires a certain block height to be reached before a certain action is taken.

Here are the steps of the scenario:

- **Given** a block with the following header:

| key | value |
| --- | ----- |
| Height | 101 |

- **Given** the query:

```  prolog
block_height(Height),
Height > 100.
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 101
gas_used: 3613
answer:
  has_more: false
  variables: ["Height"]
  results:
  - substitutions:
    - variable: Height
      expression: "101"
```
