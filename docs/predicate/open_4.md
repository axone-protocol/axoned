---
sidebar_position: 16
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

- SourceSink is an atom representing the source or sink of the stream, which is typically a URI.
- Mode is an atom representing the mode of the stream to be opened. It can be one of "read", "write", or "append".
- Stream is the stream to be opened.
- Options is a list of options. No options are currently defined, so the list should be empty.

open/4 gives True when SourceSink can be opened in Mode with the given Options.

## Virtual File System \\\(VFS\\\)

The logical module interprets on\-chain Prolog programs, relying on a Virtual Machine that isolates execution from the external environment. Consequently, the open/4 predicate doesn't access the physical file system as one might expect. Instead, it operates with a Virtual File System \(VFS\), a conceptual layer that abstracts the file system. This abstraction offers a unified view across various storage systems, adhering to the constraints imposed by blockchain technology.

This VFS extends the file concept to resources, which are identified by a Uniform Resource Identifier \(URI\). A URI specifies the access protocol for the resource, its path, and any necessary parameters.

## CosmWasm URI

The cosmwasm URI enables interaction with instantiated CosmWasm smart contract on the blockchain. The URI is used to query the smart contract and retrieve the response. The query is executed on the smart contract, and the response is returned as a stream. Query parameters are passed as part of the URI to customize the interaction with the smart contract.

Its format is as follows:

```text
cosmwasm:{contract_name}:{contract_address}?query={contract_query}[&base64Decode={true|false}]
```

where:

- \{contract\_name\}: For informational purposes, indicates the name or type of the smart contract \(e.g., "axone\-objectarium"\).
- \{contract\_address\}: Specifies the smart contract instance to query.
- \{contract\_query\}: The query to be executed on the smart contract. It is a JSON object that specifies the query payload.
- base64Decode: \(Optional\) If true, the response is base64\-decoded. Otherwise, the response is returned as is.

## Examples

### Open a resource for reading

This scenario showcases the procedure for accessing a resource stored within a CosmWasm smart contract for reading
purposes and obtaining the stream's properties.

Assuming the existence of a CosmWasm smart contract configured to store resources, we construct a URI to specifically
identify the smart contract and pinpoint the resource we aim to retrieve via a query message.

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

- **Given** the program:

```  prolog
atomic_list_concat([], '').
atomic_list_concat([H|T], Atom) :-
  atomic_list_concat(T, TAtom),
  atom_concat(H, TAtom, Atom).

resource_uri(Contract, Query, URI) :-
  uri_encoded(query_value, Query, EncodedQuery),
  atomic_list_concat(['cosmwasm:storage:', Contract, '?query=', EncodedQuery, '&base64Decode=false'], URI).
```

- **Given** the query:

```  prolog
resource_uri(
  'axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk',
  '{"object_data":{"id": "4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05"}}',
  URI),
open(URI, read, _, []).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4146
answer:
  has_more: false
  variables: ["URI"]
  results:
  - substitutions:
    - variable: URI
      expression: "'cosmwasm:storage:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D&base64Decode=false'"
```

### Open an existing resource and read its content

This scenario shows a more complex example of how to open an existing resource stored in a CosmWasm smart contract
and read its content.

The resource is opened for reading, and the content is read into a list of characters. Finally, the stream is closed.

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

- **Given** the program:

```  prolog
read_resource(Resource, Chars) :-
  open(Resource, read, Stream, []),
  read_string(Stream, _, Chars),
  close(Stream).
```

- **Given** the query:

```  prolog
read_resource('cosmwasm:storage:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D&base64Decode=false', Chars).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4142
answer:
  has_more: false
  variables: ["Chars"]
  results:
  - substitutions:
    - variable: Chars
      expression: "'Hello, World!'"
```

### Open an existing resource and read its content (base64-encoded)

This scenario is a variation of the previous one. The difference is that the smart contract returns a base64-encoded
response. For this reason, we set the `base64Decode` parameter to `true` in the query (the default value is `false`).

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
  "SGVsbG8sIFdvcmxkIQ=="
```

- **Given** the program:

```  prolog
read_resource(Resource, Chars) :-
  open(Resource, read, Stream, []),
  read_string(Stream, _, Chars),
  close(Stream).
```

- **Given** the query:

```  prolog
read_resource('cosmwasm:storage:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D&base64Decode=true', Chars).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4142
answer:
  has_more: false
  variables: ["Chars"]
  results:
  - substitutions:
    - variable: Chars
      expression: "'Hello, World!'"
```

### Try to open a non-existing resource

This scenario demonstrates the system's response to trying to open a non-existing resource.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
open('cosmwasm:storage:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=foo', read, Stream, []).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4140
answer:
  has_more: false
  variables: ["Stream"]
  results:
  - error: "error(existence_error(source_sink,cosmwasm:storage:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=foo),open/4)"
```

### Try to open a resource for writing

This scenario demonstrates the system's response to opening a resource for writing, but the resource does not allow
writing. This is the case for resources hosted in smart contracts which are read-only.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
open('cosmwasm:storage:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=foo', write, Stream, []).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4140
answer:
  has_more: false
  variables: ["Stream"]
  results:
  - error: "error(permission_error(input,stream,cosmwasm:storage:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=foo),open/4)"
```

### Try to open a resource for appending

This scenario demonstrates the system's response to opening a resource for appending, but the resource does not allow
appending. This is the case for resources hosted in smart contracts which are read-only.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
open('cosmwasm:storage:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=foo', write, Stream, []).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4140
answer:
  has_more: false
  variables: ["Stream"]
  results:
  - error: "error(permission_error(input,stream,cosmwasm:storage:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=foo),open/4)"
```

### Pass incorrect options to open/4

This scenario demonstrates the system's response to opening a resource with incorrect options.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
open('cosmwasm:storage:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=foo', read, Stream, [non_existing_option]).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4140
answer:
  has_more: false
  variables: ["Stream"]
  results:
  - error: "error(domain_error(empty_list,[non_existing_option]),open/4)"
```
