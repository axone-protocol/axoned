---
sidebar_position: 1
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# wasm_query/3

## Module

This predicate is provided by `wasm.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/wasm.pl').
```

## Description

Executes a CosmWasm smart query against the contract at Address.

- Address must be a valid Bech32 account address.
- RequestBytes is the exact query payload as bytes (typically UTF-8 JSON).
- ResponseBytes is unified with the raw response bytes returned by the contract.

Both RequestBytes and ResponseBytes use lists of integers in [0,255].

## Signature

```text
wasm_query(+Address, +RequestBytes, -ResponseBytes) is det
```

## Examples

### Query a smart contract and read raw response bytes

This scenario demonstrates how to send a smart-query payload as bytes and get the raw response bytes back.

Here are the steps of the scenario:

- **Given** the CosmWasm smart contract "axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk" and the behavior:

```  yaml
message: |
  {"ping":"pong"}
response: |
  {"ok":true}
```

- **Given** the program:

```  prolog
:- consult('/v1/lib/wasm.pl').
```

- **Given** the query:

```  prolog
wasm_query(
  'axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk',
  [123,34,112,105,110,103,34,58,34,112,111,110,103,34,125],
  ResponseBytes
).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 16910
answer:
  has_more: false
  variables: ["ResponseBytes"]
  results:
  - substitutions:
    - variable: ResponseBytes
      expression: "[123,34,111,107,34,58,116,114,117,101,125]"
```

### Reject an invalid contract address

This scenario demonstrates how wasm_query/3 validates the contract address format before calling the chain.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- consult('/v1/lib/wasm.pl').
```

- **Given** the query:

```  prolog
wasm_query('invalid-address', [123,125], ResponseBytes).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 9285
answer:
  has_more: false
  variables: ["ResponseBytes"]
  results:
  - error: "error(domain_error(valid_encoding(bech32),invalid-address),wasm_query/3)"
```

### Reject a non-byte request payload

This scenario demonstrates how wasm_query/3 raises a type error when the payload list contains values outside the byte range [0,255].

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- consult('/v1/lib/wasm.pl').
```

- **Given** the query:

```  prolog
wasm_query(
  'axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk',
  [256],
  ResponseBytes
).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 11690
answer:
  has_more: false
  variables: ["ResponseBytes"]
  results:
  - error: "error(type_error(byte,256),wasm_query/3)"
```

### Surface contract query execution failures

This scenario demonstrates that a contract query failure is surfaced as a system error in wasm_query/3.

Here are the steps of the scenario:

- **Given** the CosmWasm smart contract "axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk" and the behavior:

```  yaml
message: |
  {"ping":"pong"}
error: wasm contract execution failed
```

- **Given** the program:

```  prolog
:- consult('/v1/lib/wasm.pl').
```

- **Given** the query:

```  prolog
wasm_query(
  'axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk',
  [123,34,112,105,110,103,34,58,34,112,111,110,103,34,125],
  ResponseBytes
).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 15866
answer:
  has_more: false
  variables: ["ResponseBytes"]
  results:
  - error: "error(system_error,wasm_query/3)"
```
