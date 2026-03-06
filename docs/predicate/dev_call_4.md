---
sidebar_position: 48
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# dev_call/4

## Module

This predicate is provided by `dev.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/dev.pl').
```

## Description

Executes a transactional device call following a half-duplex protocol.

## Overview

A device is a special type of file in the virtual filesystem that implements
a transactional request-response protocol. Unlike regular files, devices follow
a strict half-duplex communication pattern with three distinct phases:

1. **Request phase**: Write operations accumulate bytes into a request buffer
2. **Commit phase**: First read operation commits the accumulated request
3. **Response phase**: Subsequent reads stream the response until EOF

Once committed (after the first read), the device transitions to read-only mode
and rejects any further write attempts with a permission error.

## Protocol Flow

```
1. open device stream in read_write mode
2. run WriteGoal(Stream)  ← builds request
3. run ReadGoal(Stream)   ← commits & reads response
4. close stream
```

The commit operation is device-specific and executes the actual transaction
(e.g., a smart contract query, a database call, etc.). Most devices require
at least one write before the first read; reading without writing typically
fails with an `invalid_request` error.

## Arguments

- `Path`: atom representing the device path in the VFS (e.g., '/v1/dev/wasm/...')
- `Type`: stream type, either `text` or `binary`
- `WriteGoal`: callable that receives Stream as first argument to build the request
- `ReadGoal`: callable that receives Stream as first argument to read the response

## Usage Notes

**⚠️ Advanced Feature**: This predicate provides low-level access to transactional
devices in the virtual filesystem of the Prolog VM.
For most use cases, prefer higher-level predicates like `wasm_query/3` for smart
contracts, which provide simpler interfaces.

Use `dev_call/4` only when you need:

- Direct control over the device protocol
- Custom request/response handling
- Integration with devices that don't have specialized predicates

Transactional devices include:

- Codecs and transforms (bech32, base64, etc.) exposed as devices
- WASM smart contract interactions (prefer `wasm_query/3` for common cases)
- Other transactional operations as needed

## Signature

```text
dev_call(+Path, +Type, :WriteGoal, :ReadGoal) is det
```

## Examples

### Echo roundtrip using dev_call/4 with meta goals

This scenario demonstrates the typical successful usage of dev_call/4:

- The write goal sends request bytes to the device
- The read goal commits the request and reads the response
- The device echoes back the exact bytes sent

Here are the steps of the scenario:

- **Given** the query:

```  prolog
echo([0,1,2,255], Echoed).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 6270
answer:
  has_more: false
  variables: ["Echoed"]
  results:
  - substitutions:
    - variable: Echoed
      expression: "[0,1,2,255]"
```

### Reading without writing fails with invalid_request

This scenario illustrates the half-duplex protocol requirement:
the device expects at least one write before the first read commits.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- consult('/v1/lib/dev.pl').

read_without_write(Result) :-
  dev_call('/v1/dev/echo', binary, no_write, read_bytes(Result)).

no_write(_).

read_bytes(Stream, Bytes) :-
  dev_read_bytes(Stream, Bytes).
```

- **Given** the query:

```  prolog
catch(
  read_without_write(_),
  error(system_error, _),
  true
).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 5021
answer:
  has_more: false
  results:
  - {}
```

### Writing after reading fails with permission error

This scenario demonstrates that once the first read commits the request,
the device transitions to read-only mode and rejects further writes.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- consult('/v1/lib/dev.pl').

write_after_read(Echoed) :-
  dev_call('/v1/dev/echo', binary, write_then_read(Echoed), no_read).

no_read(_).

write_then_read(Stream, Echoed) :-
  dev_write_bytes(Stream, [1,2,3]),
  dev_read_bytes(Stream, Echoed),
  dev_write_bytes(Stream, [4,5,6]).
```

- **Given** the query:

```  prolog
catch(
  write_after_read(_),
  error(system_error, _),
  true
).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 6415
answer:
  has_more: false
  results:
  - {}
```

### Multiple reads stream the response progressively

This scenario shows that after commit, multiple read operations
can progressively consume the response stream until EOF.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- consult('/v1/lib/dev.pl').

echo_partial(Result) :-
  dev_call('/v1/dev/echo', binary, write_all, read_partial(Result)).

write_all(Stream) :-
  dev_write_bytes(Stream, [65,66,67,68]).

read_partial(Stream, Result) :-
  get_byte(Stream, B1),
  get_byte(Stream, B2),
  Result = [B1, B2].
```

- **Given** the query:

```  prolog
echo_partial(Result).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 5803
answer:
  has_more: false
  variables: ["Result"]
  results:
  - substitutions:
    - variable: Result
      expression: "[65,66]"
```
