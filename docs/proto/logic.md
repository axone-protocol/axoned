[//]: # (This file is auto-generated. Please do not modify it yourself.)

# Protobuf Documentation

<a name="top"></a>

## 📝 Description

Lorem impsus dolor sit amet

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
