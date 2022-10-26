<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [logic/v1beta/params.proto](#logic/v1beta/params.proto)
    - [Params](#logic.v1beta.Params)
  
- [logic/v1beta/genesis.proto](#logic/v1beta/genesis.proto)
    - [GenesisState](#logic.v1beta.GenesisState)
  
- [logic/v1beta/query.proto](#logic/v1beta/query.proto)
    - [QueryServiceParamsRequest](#logic.v1beta.QueryServiceParamsRequest)
    - [QueryServiceParamsResponse](#logic.v1beta.QueryServiceParamsResponse)
  
    - [QueryService](#logic.v1beta.QueryService)
  
- [logic/v1beta/tx.proto](#logic/v1beta/tx.proto)
    - [MsgService](#logic.v1beta.MsgService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="logic/v1beta/params.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## logic/v1beta/params.proto



<a name="logic.v1beta.Params"></a>

### Params
Params defines the parameters for the module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `foo` | [string](#string) |  | foo represents a metasyntactic variable for testing purposes. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="logic/v1beta/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## logic/v1beta/genesis.proto



<a name="logic.v1beta.GenesisState"></a>

### GenesisState
GenesisState defines the logic module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#logic.v1beta.Params) |  | The state parameters for the logic module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="logic/v1beta/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## logic/v1beta/query.proto



<a name="logic.v1beta.QueryServiceParamsRequest"></a>

### QueryServiceParamsRequest
QueryServiceParamsRequest is request type for the QueryService/Params RPC method.






<a name="logic.v1beta.QueryServiceParamsResponse"></a>

### QueryServiceParamsResponse
QueryServiceParamsResponse is response type for the QueryService/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#logic.v1beta.Params) |  | params holds all the parameters of this module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="logic.v1beta.QueryService"></a>

### QueryService
QueryService defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryServiceParamsRequest](#logic.v1beta.QueryServiceParamsRequest) | [QueryServiceParamsResponse](#logic.v1beta.QueryServiceParamsResponse) | Parameters queries the parameters of the module. | GET|/okp4/okp4d/logic/params|

 <!-- end services -->



<a name="logic/v1beta/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## logic/v1beta/tx.proto


 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="logic.v1beta.MsgService"></a>

### MsgService
Msg defines the Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |

 <!-- end services -->



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

