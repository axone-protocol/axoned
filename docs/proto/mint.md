[//]: # (This file is auto-generated. Please do not modify it yourself.)

# Protobuf Documentation

<a name="top"></a>

## Purpose

The Minting Module is a fundamental component within the blockchain ecosystem, responsible for managing the dynamic
issuance of tokens to validators, thereby ensuring the stability and sustainability of the network.
This module operates on a continuous basis, recalculating parameters with each new block to maintain an optimal
token supply.

## 🧮 Calculation

The formula for calculating the yearly staking reward target (srty) is as follows:

```
srty = a • i • (c - b / bt)
```

- `srty`: Represents the yearly staking reward target.
- `a`: Denotes the total token supply within the blockchain.
- `i`: Refers to the inflation coefficient, a configurable parameter set at 0.073 by default.
- `c`: Signifies the bonding ratio adjustment coefficient.
- `b`: A dynamic variable ranging from 0 to 1, symbolizing the current bonding ratio at the current block.
- `bt`: Represents the target bonding ratio, conventionally set to 0.66.

### Yearly Inflation Rate

Upon the determination of srty, the Minting Module proceeds to convert it into a yearly inflation rate
expressed as a percentage (Iy):

```
Iy = (srty • 100) / a = (a • i • (c - b / bt) • 100) / a = i • (c - b / bt) • 100
```

- `Iy` (in percentage): Represents the yearly inflation rate.

### Block-Level Token Minting

With the yearly inflation rate in hand, the Minting Module can calculate the number of tokens to be minted for
each individual block by knowing the number of block in one years.

## Table of Contents

- [mint/v1beta1/mint.proto](#mint/v1beta1/mint.proto)
  - [Minter](#mint.v1beta1.Minter)
  - [Params](#mint.v1beta1.Params)
  
- [mint/v1beta1/genesis.proto](#mint/v1beta1/genesis.proto)
  - [GenesisState](#mint.v1beta1.GenesisState)
  
- [mint/v1beta1/query.proto](#mint/v1beta1/query.proto)
  - [QueryAnnualProvisionsRequest](#mint.v1beta1.QueryAnnualProvisionsRequest)
  - [QueryAnnualProvisionsResponse](#mint.v1beta1.QueryAnnualProvisionsResponse)
  - [QueryInflationRequest](#mint.v1beta1.QueryInflationRequest)
  - [QueryInflationResponse](#mint.v1beta1.QueryInflationResponse)
  - [QueryParamsRequest](#mint.v1beta1.QueryParamsRequest)
  - [QueryParamsResponse](#mint.v1beta1.QueryParamsResponse)
  
  - [Query](#mint.v1beta1.Query)
  
- [mint/v1beta1/tx.proto](#mint/v1beta1/tx.proto)
  - [MsgUpdateParams](#mint.v1beta1.MsgUpdateParams)
  - [MsgUpdateParamsResponse](#mint.v1beta1.MsgUpdateParamsResponse)
  
  - [Msg](#mint.v1beta1.Msg)
  
- [Scalar Value Types](#scalar-value-types)

<a name="mint/v1beta1/mint.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## mint/v1beta1/mint.proto

<a name="mint.v1beta1.Minter"></a>

### Minter

Minter represents the minting state within the blockchain, tasked with the continuous calculation and distribution
of tokens to validators.

This calculation occurs with each block, where the minting module dynamically recalculates the annual `inflation` rate.
Using the resulting inflation rate, it deduces the quantity of tokens to be provisioned for the upcoming year.

Furthermore, based on the current block's position within the year, it computes the exact number of tokens to be
minted for that specific block.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `inflation` | [string](#string) |  | Represents the current annual inflation rate. |
| `annual_provisions` | [string](#string) |  | Represents the current annual expected provisions. |

<a name="mint.v1beta1.Params"></a>

### Params

Params holds parameters for the mint module.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `mint_denom` | [string](#string) |  | Type of coin to mint |
| `inflation_coef` | [string](#string) |  | Annual inflation coefficient |
| `bonding_adjustment` | [string](#string) |  | Bonding adjustment |
| `target_bonding_ratio` | [string](#string) |  | Represent the target bonding ratio to reach |
| `blocks_per_year` | [uint64](#uint64) |  | Estimated blocks per year |

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

 [//]: # (end services)

<a name="mint/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## mint/v1beta1/genesis.proto

<a name="mint.v1beta1.GenesisState"></a>

### GenesisState

GenesisState defines the mint module's genesis state.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `minter` | [Minter](#mint.v1beta1.Minter) |  | minter is a space for holding current inflation information. |
| `params` | [Params](#mint.v1beta1.Params) |  | params defines all the parameters of the module. |

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

 [//]: # (end services)

<a name="mint/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## mint/v1beta1/query.proto

<a name="mint.v1beta1.QueryAnnualProvisionsRequest"></a>

### QueryAnnualProvisionsRequest

QueryAnnualProvisionsRequest is the request type for the
Query/AnnualProvisions RPC method.

<a name="mint.v1beta1.QueryAnnualProvisionsResponse"></a>

### QueryAnnualProvisionsResponse

QueryAnnualProvisionsResponse is the response type for the
Query/AnnualProvisions RPC method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `annual_provisions` | [bytes](#bytes) |  | annual_provisions is the current minting annual provisions value. |

<a name="mint.v1beta1.QueryInflationRequest"></a>

### QueryInflationRequest

QueryInflationRequest is the request type for the Query/Inflation RPC method.

<a name="mint.v1beta1.QueryInflationResponse"></a>

### QueryInflationResponse

QueryInflationResponse is the response type for the Query/Inflation RPC
method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `inflation` | [bytes](#bytes) |  | inflation is the current minting inflation value. |

<a name="mint.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest

QueryParamsRequest is the request type for the Query/Params RPC method.

<a name="mint.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse

QueryParamsResponse is the response type for the Query/Params RPC method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#mint.v1beta1.Params) |  | params defines the parameters of the module. |

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

<a name="mint.v1beta1.Query"></a>

### Query

Query provides defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#mint.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#mint.v1beta1.QueryParamsResponse) | Params returns the total set of minting parameters. | GET|/cosmos/mint/v1beta1/params|
| `Inflation` | [QueryInflationRequest](#mint.v1beta1.QueryInflationRequest) | [QueryInflationResponse](#mint.v1beta1.QueryInflationResponse) | Inflation returns the current minting inflation value. | GET|/cosmos/mint/v1beta1/inflation|
| `AnnualProvisions` | [QueryAnnualProvisionsRequest](#mint.v1beta1.QueryAnnualProvisionsRequest) | [QueryAnnualProvisionsResponse](#mint.v1beta1.QueryAnnualProvisionsResponse) | AnnualProvisions current minting annual provisions value. | GET|/cosmos/mint/v1beta1/annual_provisions|

 [//]: # (end services)

<a name="mint/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## mint/v1beta1/tx.proto

<a name="mint.v1beta1.MsgUpdateParams"></a>

### MsgUpdateParams

MsgUpdateParams is the Msg/UpdateParams request type.

Since: cosmos-sdk 0.47

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  | authority is the address that controls the module (defaults to x/gov unless overwritten). |
| `params` | [Params](#mint.v1beta1.Params) |  | params defines the x/mint parameters to update.

NOTE: All parameters must be supplied. |

<a name="mint.v1beta1.MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse

MsgUpdateParamsResponse defines the response structure for executing a
MsgUpdateParams message.

Since: cosmos-sdk 0.47

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

<a name="mint.v1beta1.Msg"></a>

### Msg

Msg defines the x/mint Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `UpdateParams` | [MsgUpdateParams](#mint.v1beta1.MsgUpdateParams) | [MsgUpdateParamsResponse](#mint.v1beta1.MsgUpdateParamsResponse) | UpdateParams defines a governance operation for updating the x/mint module parameters. The authority is defaults to the x/gov module account.

Since: cosmos-sdk 0.47 | |

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
