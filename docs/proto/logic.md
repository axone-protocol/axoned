[//]: # (This file is auto-generated. Please do not modify it yourself.)

# Protobuf Documentation

<a name="top"></a>

## 📝 Description

This module implements a [Prolog](https://en.wikipedia.org/wiki/Prolog) logic interpreter (go-native) and
its [unification](https://en.wikipedia.org/wiki/Unification_(computer_science))
algorithm to evaluate logical expressions against the current state of the blockchain.

This distinctive module allows for the creation of advanced, goal-oriented queries and logical systems that can be
applied
to a wide range of use cases, while still maintaining the determinism and predictability of blockchain technology. It
also features a collection of predefined, blockchain-specific predicates that can be used to access information about
the state of the blockchain.

## Concepts

### Program

A program is a text that is parsed and compiled by the interpreter. A program is composed of a set of predicates, which
are defined by the user and can be used to express the desired query logic.

#### Predicate

A predicate is a statement that describes a relationship between one or more variables or constants. A predicate
consists of a name followed by zero or more arguments.

#### Rule & Fact

A rule is a statement that describes a relationship between one or more variables or constants, similar to a predicate.
However, unlike a predicate, a rule also specifies one or more conditions that must be true in order for the
relationship described by the rule to hold.

A rule has the following format:

```prolog
head :- body.
```

The symbol `:-` is called the "if-then" operator, and it means that the relationship described in the head of the rule
holds only if the conditions in the body are true.

For example:

```prolog
grandfather(X,Y) :- father(X,Z), father(Z,Y). # X is the grandfather of Y if X is the father of Z and Z is the father of Y.
```

A fact is a special type of rule that has no body (with no `:-` and no conditions). A fact has the following format:

```prolog
head.
```

For instance:

```prolog
father(john, mary). # john is the father of mary.
```

#### Variable

A variable is a predicate argument that is used as a placeholder for a value. It can represent any type of data, such as
numbers, strings, or lists.

Variables are denoted by a name that starts with an uppercase letter, for example `X` or `Foo`.

For instance:

```prolog
father(X, mary). # ask for all X that are the father of mary.
```

### Query

A query is a statement used to retrieve information from the blockchain. It can be sent against a program, but this
is optional. The interpreter evaluates the query and returns the result to the caller. Queries can be submitted to a
module using the `Ask` message.

#### `Ask`

The `Ask` message is used to submit a query to the module. It has the following format:

```text
{
  string Program
  string Query
}
```

The `Program` field is optional. If it is not specified, the query is just evaluated against the current state of the
blockchain.
If it is specified, the query is evaluated against the program that is passed as an argument.

For instance:

```text
{
  Program: "father(john, mary)."
  Query: "father(X, mary)."
}
```

Gives:

```json
{
  "height": "7235",
  "gas_used": "9085",
  "answer": {
    "has_more": false,
    "variables": [
      "X"
    ],
    "results": [
      {
        "substitutions": [
          {
            "variable": "X",
            "expression": "john"
          }
        ]
      }
    ]
  }
}
```

The logic module supports chain-specific predicates that can be used to query the state of the blockchain. For example,
the `chain_id` predicate can be used to retrieve the chain ID of the current blockchain. Several other predicates are
available, such as `block_height`, `block_time`... Please refer to the go documentation for the full list of available
predicates.

For instance:

```prolog
chain_id(X). # ask for the chain ID.
```

#### Response

The response is an object that contains the following fields:

- `height`: the height of the block at which the query was evaluated.
- `gas_used`: the amount of gas used to evaluate the query.
- `answer`: the result of the query. It is an object that contains the following fields:
  - `has_more`: a boolean that indicates whether there are more results to be retrieved. It's just informative since no
    more results can be retrieved.
  - `variables`: an array of strings that contains the names of the variables that were used in the query.
  - `results`: an array of objects that contains the solutions of the query. Each result is an object that contains the
    following fields:
    - `error`: an optional string that contains an error message if the query failed for the current solution.
    - `substitutions`: an array of objects that contains the substitutions that were made to satisfy the query. A
      substitution is a set of variable-value pairs that is used to replace variables with constants. A substitution
      is the result of unification. A substitution is used to replace variables with constants when evaluating a rule.

## Performance

The performance of the logic module is closely tied to the complexity of the query and the size of the program. To
optimize performance, especially in a constrained environment like the blockchain, it is important to minimize the size of the
program.
Keep in mind that he module uses [backtracking](https://en.wikipedia.org/wiki/Backtracking) to search for solutions,
making it most effective when used for queries that are satisfiable. Indeed, if the query is not satisfiable, the module will
attempt to find a solution by [backtracking](https://en.wikipedia.org/wiki/Backtracking) and searching through possible
solutions for an extended period before ultimately being canceled.

## Gas

The `Ask` message incurs gas consumption, which is calculated as the sum of the gas used to evaluate each predicate during
the query evaluation process. Each predicate has a fixed gas cost that is based on its complexity.

While querying the module does not require any fees, the use of gas serves as a mechanism to limit the size and
complexity of the query, ensuring optimal performance and fairness.

## Security

The logic module is a deterministic program that is executed in a sandboxed environment and does not have the ability
to submit transactions or make changes to the blockchain's state. It is therefore safe to use.

To control the cpu and memory usage of the module, the module is limited by several different mechanisms:

- `max_size`: the maximum size of the program that can be evaluated.
- `max_result_count`: the maximum number of results that can be returned by a query.

The existing `query-gas-limit` configuration present in the `app.toml` can be used to constraint gas usage when not used
in the context of a transaction.

Additional limitations are being considered for the future, such as restricting the number of variables that can be
utilized within a query, or limiting the depth of the backtracking algorithm.

## Table of Contents

- [logic/v1beta3/params.proto](#logic/v1beta3/params.proto)
  - [GasPolicy](#logic.v1beta3.GasPolicy)
  - [Limits](#logic.v1beta3.Limits)
  - [Params](#logic.v1beta3.Params)
  
- [logic/v1beta3/types.proto](#logic/v1beta3/types.proto)
  - [Answer](#logic.v1beta3.Answer)
  - [ProgramPublication](#logic.v1beta3.ProgramPublication)
  - [Result](#logic.v1beta3.Result)
  - [StoredProgram](#logic.v1beta3.StoredProgram)
  - [Substitution](#logic.v1beta3.Substitution)
  
- [logic/v1beta3/genesis.proto](#logic/v1beta3/genesis.proto)
  - [GenesisProgramPublication](#logic.v1beta3.GenesisProgramPublication)
  - [GenesisState](#logic.v1beta3.GenesisState)
  - [GenesisStoredProgram](#logic.v1beta3.GenesisStoredProgram)
  
- [logic/v1beta3/query.proto](#logic/v1beta3/query.proto)
  - [QueryServiceAskRequest](#logic.v1beta3.QueryServiceAskRequest)
  - [QueryServiceAskResponse](#logic.v1beta3.QueryServiceAskResponse)
  - [QueryServiceParamsRequest](#logic.v1beta3.QueryServiceParamsRequest)
  - [QueryServiceParamsResponse](#logic.v1beta3.QueryServiceParamsResponse)
  
  - [QueryService](#logic.v1beta3.QueryService)
  
- [logic/v1beta3/tx.proto](#logic/v1beta3/tx.proto)
  - [MsgStoreProgram](#logic.v1beta3.MsgStoreProgram)
  - [MsgStoreProgramResponse](#logic.v1beta3.MsgStoreProgramResponse)
  - [MsgUpdateParams](#logic.v1beta3.MsgUpdateParams)
  - [MsgUpdateParamsResponse](#logic.v1beta3.MsgUpdateParamsResponse)
  
  - [MsgService](#logic.v1beta3.MsgService)
  
- [Scalar Value Types](#scalar-value-types)

<a name="logic/v1beta3/params.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## logic/v1beta3/params.proto

<a name="logic.v1beta3.GasPolicy"></a>

### GasPolicy

GasPolicy defines the coefficients used to translate VM metering units into SDK gas.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `compute_coeff` | [uint64](#uint64) |  | compute_coeff applies to Instruction, ArithNode, and CompareStep VM meter kinds. If set to 0, the value considered is 1. |
| `memory_coeff` | [uint64](#uint64) |  | memory_coeff applies to CopyNode and ListCell VM meter kinds. If set to 0, the value considered is 1. |
| `unify_coeff` | [uint64](#uint64) |  | unify_coeff applies to UnifyStep VM meter kind. If set to 0, the value considered is 1. |
| `io_coeff` | [uint64](#uint64) |  | io_coeff applies to the total size in bytes of the user-supplied program and query sources, and to codec device I/O buffered by the VFS. If set to 0, the value considered is 1. |

<a name="logic.v1beta3.Limits"></a>

### Limits

Limits defines the limits of the logic module.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_size` | [uint64](#uint64) |  | max_size specifies the maximum total size, in bytes, accepted for the user-supplied program and query sources. A value of 0 means that there is no limit on the total source size. |
| `max_result_count` | [uint64](#uint64) |  | max_result_count specifies the maximum number of results that can be requested for a query. A value of 0 means that there is no limit on the number of results. |
| `max_user_output_size` | [uint64](#uint64) |  | max_user_output_size specifies the maximum number of bytes to keep in the user output. If the user output exceeds this size, the interpreter will overwrite the oldest bytes with the new ones to keep the size constant. A value of 0 means the user output is disabled. |
| `max_variables` | [uint64](#uint64) |  | max_variables specifies the maximum number of variables that can be create by the interpreter. A value of 0 means that there is no limit on the number of variables. |

<a name="logic.v1beta3.Params"></a>

### Params

Params defines all the configuration parameters of the "logic" module.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `limits` | [Limits](#logic.v1beta3.Limits) |  | Limits defines the limits of the logic module. The limits are used to prevent the interpreter from running for too long. If the interpreter runs for too long, the execution will be aborted. |
| `gas_policy` | [GasPolicy](#logic.v1beta3.GasPolicy) |  | GasPolicy defines the coefficients used to translate VM metering units into SDK gas. |

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

 [//]: # (end services)

<a name="logic/v1beta3/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## logic/v1beta3/types.proto

<a name="logic.v1beta3.Answer"></a>

### Answer

Answer represents the answer to a logic query.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `has_more` | [bool](#bool) |  | has_more specifies if there are more solutions than the ones returned. |
| `variables` | [string](#string) | repeated | variables represent all the variables in the query. |
| `results` | [Result](#logic.v1beta3.Result) | repeated | results represent all the results of the query. |

<a name="logic.v1beta3.ProgramPublication"></a>

### ProgramPublication

ProgramPublication represents the publication metadata of a program by a publisher.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `published_at` | [int64](#int64) |  | published_at is the block timestamp (Unix seconds) of publication for this publisher. |

<a name="logic.v1beta3.Result"></a>

### Result

Result represents the result of a query.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `error` | [string](#string) |  | error specifies the error message if the query caused an error. |
| `substitutions` | [Substitution](#logic.v1beta3.Substitution) | repeated | substitutions represent all the substitutions made to the variables in the query to obtain the answer. |

<a name="logic.v1beta3.StoredProgram"></a>

### StoredProgram

StoredProgram represents a program source with its storage metadata.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `source` | [string](#string) |  | source is the original Prolog source. |
| `created_at` | [int64](#int64) |  | created_at is the block timestamp (Unix seconds) of artifact creation. |
| `source_size` | [uint64](#uint64) |  | source_size is the source size in bytes. |

<a name="logic.v1beta3.Substitution"></a>

### Substitution

Substitution represents a substitution made to the variables in the query to obtain the answer.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `variable` | [string](#string) |  | variable is the name of the variable. |
| `expression` | [string](#string) |  | expression is the value substituted for the variable, represented directly as a Prolog term (e.g., atom, number, compound). |

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

 [//]: # (end services)

<a name="logic/v1beta3/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## logic/v1beta3/genesis.proto

<a name="logic.v1beta3.GenesisProgramPublication"></a>

### GenesisProgramPublication

GenesisProgramPublication associates a publisher and program_id with publication metadata.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `publisher` | [string](#string) |  | publisher is the bech32 account address that published the program. |
| `program_id` | [string](#string) |  | program_id is the SHA-256 hash of the published program source encoded as lowercase hexadecimal. |
| `publication` | [ProgramPublication](#logic.v1beta3.ProgramPublication) |  | publication is the publication metadata for this publisher/program pair. |

<a name="logic.v1beta3.GenesisState"></a>

### GenesisState

GenesisState defines the logic module's genesis state.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#logic.v1beta3.Params) |  | The state parameters for the logic module. |
| `stored_programs` | [GenesisStoredProgram](#logic.v1beta3.GenesisStoredProgram) | repeated | stored_programs are the canonical immutable program artifacts keyed by program_id. |
| `program_publications` | [GenesisProgramPublication](#logic.v1beta3.GenesisProgramPublication) | repeated | program_publications are the user-scoped immutable publication views pointing to stored artifacts. |

<a name="logic.v1beta3.GenesisStoredProgram"></a>

### GenesisStoredProgram

GenesisStoredProgram associates a program_id with its canonical stored artifact.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program_id` | [string](#string) |  | program_id is the SHA-256 hash of the program source encoded as lowercase hexadecimal. |
| `program` | [StoredProgram](#logic.v1beta3.StoredProgram) |  | program is the canonical immutable stored program artifact. |

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

 [//]: # (end services)

<a name="logic/v1beta3/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## logic/v1beta3/query.proto

<a name="logic.v1beta3.QueryServiceAskRequest"></a>

### QueryServiceAskRequest

QueryServiceAskRequest is request type for the QueryService/Ask RPC method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program` | [string](#string) |  | program is the logic program to be queried. |
| `query` | [string](#string) |  | query is the query string to be executed. |
| `limit` | [uint64](#uint64) |  | limit specifies the maximum number of solutions to be returned. This field is governed by max_result_count, which defines the upper limit of results that may be requested per query. If this field is not explicitly set, a default value of 1 is applied. |

<a name="logic.v1beta3.QueryServiceAskResponse"></a>

### QueryServiceAskResponse

QueryServiceAskResponse is response type for the QueryService/Ask RPC method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [uint64](#uint64) |  | height is the block height at which the query was executed. |
| `gas_used` | [uint64](#uint64) |  | gas_used is the amount of gas used to execute the query. |
| `answer` | [Answer](#logic.v1beta3.Answer) |  | answer is the answer to the query. |
| `user_output` | [string](#string) |  | user_output is the output of the query execution, if any. the length of the output is limited by the max_query_output_size parameter. |

<a name="logic.v1beta3.QueryServiceParamsRequest"></a>

### QueryServiceParamsRequest

QueryServiceParamsRequest is request type for the QueryService/Params RPC method.

<a name="logic.v1beta3.QueryServiceParamsResponse"></a>

### QueryServiceParamsResponse

QueryServiceParamsResponse is response type for the QueryService/Params RPC method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#logic.v1beta3.Params) |  | params holds all the parameters of this module. |

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

<a name="logic.v1beta3.QueryService"></a>

### QueryService

QueryService defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryServiceParamsRequest](#logic.v1beta3.QueryServiceParamsRequest) | [QueryServiceParamsResponse](#logic.v1beta3.QueryServiceParamsResponse) | Params queries all parameters for the logic module. | GET|/axone-protocol/axoned/logic/params|
| `Ask` | [QueryServiceAskRequest](#logic.v1beta3.QueryServiceAskRequest) | [QueryServiceAskResponse](#logic.v1beta3.QueryServiceAskResponse) | Ask executes a logic query and returns the solutions found. Since the query is without any side-effect, the query is not executed in the context of a transaction and no fee is charged for this, but the execution is constrained by the current limits configured in the module. | GET|/axone-protocol/axoned/logic/ask|

 [//]: # (end services)

<a name="logic/v1beta3/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## logic/v1beta3/tx.proto

<a name="logic.v1beta3.MsgStoreProgram"></a>

### MsgStoreProgram

MsgStoreProgram defines a Msg for storing a Prolog program source as a user library.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `publisher` | [string](#string) |  | publisher is the bech32 account address publishing the program artifact. After publication, this exact address is used as the <identity> path segment in the logic module virtual file system path /v1/var/lib/logic/users/<identity>/programs/<program_id>.pl. This is the path that Prolog code can load through consult/1. |
| `source` | [string](#string) |  | source is the Prolog program source to parse and store. |

<a name="logic.v1beta3.MsgStoreProgramResponse"></a>

### MsgStoreProgramResponse

MsgStoreProgramResponse defines the response for executing a MsgStoreProgram.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program_id` | [string](#string) |  | program_id is the SHA-256 hash of the program source (lowercase hexadecimal). After publication, this exact identifier is used as the <program_id> path segment in the logic module virtual file system path /v1/var/lib/logic/users/<identity>/programs/<program_id>.pl. This is the path that Prolog code can load through consult/1. |

<a name="logic.v1beta3.MsgUpdateParams"></a>

### MsgUpdateParams

MsgUpdateParams defines a Msg for updating the x/logic module parameters.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  | authority is the address of the governance account. |
| `params` | [Params](#logic.v1beta3.Params) |  | params defines the x/logic parameters to update. NOTE: All parameters must be supplied. |

<a name="logic.v1beta3.MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse

MsgUpdateParamsResponse defines the response structure for executing a
MsgUpdateParams message.

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

<a name="logic.v1beta3.MsgService"></a>

### MsgService

MsgService defines the transaction service for the logic module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `UpdateParams` | [MsgUpdateParams](#logic.v1beta3.MsgUpdateParams) | [MsgUpdateParamsResponse](#logic.v1beta3.MsgUpdateParamsResponse) | UpdateParams defined a governance operation for updating the x/logic module parameters. The authority is hard-coded to the Cosmos SDK x/gov module account | |
| `StoreProgram` | [MsgStoreProgram](#logic.v1beta3.MsgStoreProgram) | [MsgStoreProgramResponse](#logic.v1beta3.MsgStoreProgramResponse) | StoreProgram validates a Prolog user library source and stores its canonical artifact if needed. Artifact identity is content-addressed: program_id = sha256(source). The endpoint is idempotent for the same publisher + same source, and also when different publishers submit the same source. After a successful call, the published program is exposed through the logic module virtual file system at the immutable path /v1/var/lib/logic/users/<identity>/programs/<program_id>.pl. This path is intended to be loaded from Prolog, for example with consult/1. | |

 [//]: # (end services)

## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |
