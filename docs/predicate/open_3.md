---
sidebar_position: 21
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# open/3

## Description

`open/3` is a predicate which opens a stream to a source or sink. This predicate is a shorthand for open/4 with an empty list of options.

## Signature

```text
open(+SourceSink, +Mode, -Stream)
```

where:

- SourceSink is an atom representing the source or sink of the stream, which is typically a URI.
- Mode is an atom representing the mode of the stream to be opened. It can be one of "read", "write", or "append".
- Stream is the stream to be opened.

open/3 gives True when SourceSink can be opened in Mode.

## Examples

### Open a resource for reading

This scenario showcases the procedure for accessing a resource stored within a CosmWasm smart contract for reading
purposes and obtaining the stream's properties.

See the `open/4` predicate for a more detailed example.

Here are the steps of the scenario:

- **Given** the CosmWasm smart contract "axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk" and the behavior:

```  yaml
message: |
  {
    "object_data": {
      "id": "4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05"
    }
  }
response: |
  Hello, World!
```

- **Given** the query:

```  prolog
open(
  'cosmwasm:storage:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D&base64Decode=false',
  read,
  _
).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3613
answer:
  has_more: false
  variables:
  results:
  - substitutions:
```
