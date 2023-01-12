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
    "success": true,
    "has_more": false,
    "variables": [
      "X"
    ],
    "results": [
      {
        "substitutions": [
          {
            "variable": "X",
            "term": {
              "name": "john",
              "arguments": []
            }
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
  - `success`: a boolean that indicates whether the query was successful or not. Successful means that solutions were
    found, i.e. the query was satisfiable.
  - `has_more`: a boolean that indicates whether there are more results to be retrieved. It's just informative since no
    more results can be retrieved.
  - `variables`: an array of strings that contains the names of the variables that were used in the query.
  - `results`: an array of objects that contains the results of the query. Each result is an object that contains the
    following fields:
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

- `max_gas`: the maximum amount of gas that can be used to evaluate a query.
- `max_size`: the maximum size of the program that can be evaluated.
- `max_result_count`: the maximum number of results that can be returned by a query.

Additional limitations are being considered for the future, such as restricting the number of variables that can be
utilized within a query, or limiting the depth of the backtracking algorithm.

## Table of Contents

- [logic/v1beta/params.proto](#logic/v1beta/params.proto)
  - [Interpreter](#logic.v1beta.Interpreter)
  - [Limits](#logic.v1beta.Limits)
  - [Params](#logic.v1beta.Params)
  
- [logic/v1beta/genesis.proto](#logic/v1beta/genesis.proto)
  - [GenesisState](#logic.v1beta.GenesisState)
  
- [logic/v1beta/types.proto](#logic/v1beta/types.proto)
  - [Answer](#logic.v1beta.Answer)
  - [Result](#logic.v1beta.Result)
  - [Substitution](#logic.v1beta.Substitution)
  - [Term](#logic.v1beta.Term)
  
- [logic/v1beta/query.proto](#logic/v1beta/query.proto)
  - [QueryServiceAskRequest](#logic.v1beta.QueryServiceAskRequest)
  - [QueryServiceAskResponse](#logic.v1beta.QueryServiceAskResponse)
  - [QueryServiceParamsRequest](#logic.v1beta.QueryServiceParamsRequest)
  - [QueryServiceParamsResponse](#logic.v1beta.QueryServiceParamsResponse)
  
  - [QueryService](#logic.v1beta.QueryService)
  
- [logic/v1beta/tx.proto](#logic/v1beta/tx.proto)
  - [MsgService](#logic.v1beta.MsgService)
  
- [Scalar Value Types](#scalar-value-types)

<a name="logic/v1beta/params.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## logic/v1beta/params.proto

<a name="logic.v1beta.Interpreter"></a>

### Interpreter

Interpreter defines the various parameters for the interpreter.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `registered_predicates` | [string](#string) | repeated | registered_predicates specifies the list of registered predicates/operators, in the form of: `<predicate_name>/<arity>`. For instance: `findall/3`. If not specified, the default set of predicates/operators will be registered. |
| `bootstrap` | [string](#string) |  | bootstrap specifies the initial program to run when booting the logic interpreter. If not specified, the default boot sequence will be executed. |

<a name="logic.v1beta.Limits"></a>

### Limits

Limits defines the limits of the logic module.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_gas` | [string](#string) |  | max_gas specifies the maximum amount of computing power, measured in "gas," that is allowed to be consumed when executing a request by the interpreter. The interpreter calculates the gas consumption based on the number and type of operations that are executed, as well as, in some cases, the complexity of the processed data. nil value remove max gas limitation. |
| `max_size` | [string](#string) |  | max_size specifies the maximum size, in bytes, that is accepted for a program. nil value remove size limitation. |
| `max_result_count` | [string](#string) |  | max_result_count specifies the maximum number of results that can be requested for a query. nil value remove max result count limitation. |

<a name="logic.v1beta.Params"></a>

### Params

Params defines all the configuration parameters of the "logic" module.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `interpreter` | [Interpreter](#logic.v1beta.Interpreter) |  | Interpreter specifies the parameter for the logic interpreter. |
| `limits` | [Limits](#logic.v1beta.Limits) |  | Limits defines the limits of the logic module. The limits are used to prevent the interpreter from running for too long. If the interpreter runs for too long, the execution will be aborted. |

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

 [//]: # (end services)

<a name="logic/v1beta/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## logic/v1beta/genesis.proto

<a name="logic.v1beta.GenesisState"></a>

### GenesisState

GenesisState defines the logic module's genesis state.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#logic.v1beta.Params) |  | The state parameters for the logic module. |

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

 [//]: # (end services)

<a name="logic/v1beta/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## logic/v1beta/types.proto

<a name="logic.v1beta.Answer"></a>

### Answer

Answer represents the answer to a logic query.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `success` | [bool](#bool) |  | result is the result of the query. |
| `has_more` | [bool](#bool) |  | has_more specifies if there are more solutions than the ones returned. |
| `variables` | [string](#string) | repeated | variables represent all the variables in the query. |
| `results` | [Result](#logic.v1beta.Result) | repeated | results represent all the results of the query. |

<a name="logic.v1beta.Result"></a>

### Result

Result represents the result of a query.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `substitutions` | [Substitution](#logic.v1beta.Substitution) | repeated | substitutions represent all the substitutions made to the variables in the query to obtain the answer. |

<a name="logic.v1beta.Substitution"></a>

### Substitution

Substitution represents a substitution made to the variables in the query to obtain the answer.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `variable` | [string](#string) |  | variable is the name of the variable. |
| `term` | [Term](#logic.v1beta.Term) |  | term is the term that the variable is substituted with. |

<a name="logic.v1beta.Term"></a>

### Term

Term is the representation of a piece of data and can be a constant, a variable, or an atom.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | name is the name of the term. |
| `arguments` | [Term](#logic.v1beta.Term) | repeated | arguments are the arguments of the term, which can be constants, variables, or atoms. |

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

 [//]: # (end services)

<a name="logic/v1beta/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## logic/v1beta/query.proto

<a name="logic.v1beta.QueryServiceAskRequest"></a>

### QueryServiceAskRequest

QueryServiceAskRequest is request type for the QueryService/Ask RPC method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program` | [string](#string) |  | program is the logic program to be queried. |
| `query` | [string](#string) |  | query is the query string to be executed. |

<a name="logic.v1beta.QueryServiceAskResponse"></a>

### QueryServiceAskResponse

QueryServiceAskResponse is response type for the QueryService/Ask RPC method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [uint64](#uint64) |  | height is the block height at which the query was executed. |
| `gas_used` | [uint64](#uint64) |  | gas_used is the amount of gas used to execute the query. |
| `answer` | [Answer](#logic.v1beta.Answer) |  | answer is the answer to the query. |

<a name="logic.v1beta.QueryServiceParamsRequest"></a>

### QueryServiceParamsRequest

QueryServiceParamsRequest is request type for the QueryService/Params RPC method.

<a name="logic.v1beta.QueryServiceParamsResponse"></a>

### QueryServiceParamsResponse

QueryServiceParamsResponse is response type for the QueryService/Params RPC method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#logic.v1beta.Params) |  | params holds all the parameters of this module. |

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

<a name="logic.v1beta.QueryService"></a>

### QueryService

QueryService defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryServiceParamsRequest](#logic.v1beta.QueryServiceParamsRequest) | [QueryServiceParamsResponse](#logic.v1beta.QueryServiceParamsResponse) | Params queries all parameters for the logic module. | GET|/okp4/okp4d/logic/params|
| `Ask` | [QueryServiceAskRequest](#logic.v1beta.QueryServiceAskRequest) | [QueryServiceAskResponse](#logic.v1beta.QueryServiceAskResponse) | Ask executes a logic query and returns the solutions found. Since the query is without any side-effect, the query is not executed in the context of a transaction and no fee is charged for this, but the execution is constrained by the current limits configured in the module. | GET|/okp4/okp4d/logic/ask|

 [//]: # (end services)

<a name="logic/v1beta/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## logic/v1beta/tx.proto

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

<a name="logic.v1beta.MsgService"></a>

### MsgService

MsgService defines the service for the logic module.
Do nothing for now as the service is without any side effects.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |

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
