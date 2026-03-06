---
sidebar_position: 92
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# open/4

## Description

`open/4` is a predicate which opens a stream to a source or sink.

## Signature

```text
open(+SourceSink, +Mode, -Stream, +Options)
```

where:

- SourceSink is an atom representing the source or sink of the stream in the virtual file system.
- Mode is an atom representing the mode of the stream to be opened \(for example "read", "write", "append", "read\_write"\).
- Stream is the stream to be opened.
- Options is a list of stream options.

open/4 gives True when SourceSink can be opened in Mode with the given Options.

## Virtual File System \\\(VFS\\\)

The logical module interprets on\-chain Prolog programs, relying on a Virtual Machine that isolates execution from the external environment. Consequently, the open/4 predicate doesn't access the physical file system as one might expect. Instead, it operates with a Virtual File System \(VFS\), a conceptual layer that abstracts the file system. This abstraction offers a unified view across various storage systems, adhering to the constraints imposed by blockchain technology.

This VFS extends the file concept to module\-provided resources and devices exposed as paths, for example:

- immutable snapshots under /v1/sys/\*
- transactional devices under /v1/dev/\*

## Examples

### Open a resource for reading

This scenario demonstrates how to build a VFS path and open it in read mode.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
resource_path(Path) :-
  atomic_list_concat(['/v1', 'sys', 'header', 'height'], '/', Path).
```

- **Given** the query:

```  prolog
resource_path(Path),
open(Path, read, _, []).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 6864
answer:
  has_more: false
  variables: ["Path"]
  results:
  - substitutions:
    - variable: Path
      expression: "'/v1/sys/header/height'"
```

### Open an existing resource and read its Prolog term

This scenario demonstrates how to open a text resource, read one term, and close the stream.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
read_height(Path, Height) :-
  open(Path, read, Stream, [type(text)]),
  read_term(Stream, Height, []),
  close(Stream).
```

- **Given** the query:

```  prolog
read_height('/v1/sys/header/height', Height).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4100
answer:
  has_more: false
  variables: ["Height"]
  results:
  - substitutions:
    - variable: Height
      expression: "42"
```

### Open a wasm query endpoint in read_write mode

This scenario demonstrates how to write request bytes, then read response bytes from a transactional endpoint.

Here are the steps of the scenario:

- **Given** the CosmWasm smart contract "axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk" and the behavior:

```  yaml
message: |
  {}
response: |
  {"ok":true}
```

- **Given** the program:

```  prolog
read_all_bytes(Stream, Bytes) :-
  get_byte(Stream, Byte),
  ( Byte =:= -1 ->
      Bytes = []
  ; Bytes = [Byte | Rest],
    read_all_bytes(Stream, Rest)
  ).

wasm_roundtrip(Address, ResponseBytes) :-
  atom_concat('/v1/dev/wasm/', Address, Prefix),
  atom_concat(Prefix, '/query', Path),
  open(Path, read_write, Stream, [type(binary)]),
  put_byte(Stream, 123),
  put_byte(Stream, 125),
  read_all_bytes(Stream, ResponseBytes),
  close(Stream).
```

- **Given** the query:

```  prolog
wasm_roundtrip('axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk', ResponseBytes).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 5641
answer:
  has_more: false
  variables: ["ResponseBytes"]
  results:
  - substitutions:
    - variable: ResponseBytes
      expression: "[123,34,111,107,34,58,116,114,117,101,125]"
```

### Try to open a non-existing resource

This scenario demonstrates the system's response to trying to open a non-existing resource.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
open('foo:bar', read, Stream, []).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3935
answer:
  has_more: false
  variables: ["Stream"]
  results:
  - error: "error(existence_error(source_sink,foo:bar),open/4)"
```

### Try to open a read-only resource for writing

This scenario demonstrates the system's response to opening a snapshot path in write mode.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
open('/v1/sys/header/height', write, Stream, []).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3950
answer:
  has_more: false
  variables: ["Stream"]
  results:
  - error: "error(permission_error(open,source_sink,/v1/sys/header/height),open/4)"
```

### Try to open a read-only resource for appending

This scenario demonstrates the system's response to opening a snapshot path in append mode.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
open('/v1/sys/header/height', append, Stream, []).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3951
answer:
  has_more: false
  variables: ["Stream"]
  results:
  - error: "error(permission_error(open,source_sink,/v1/sys/header/height),open/4)"
```

### Pass incorrect options to open/4

This scenario demonstrates the system's response to opening a resource with incorrect options.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
open('/v1/sys/header/height', read, Stream, [non_existing_option]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3971
answer:
  has_more: false
  variables: ["Stream"]
  results:
  - error: "error(domain_error(stream_option,non_existing_option),open/4)"
```
