[//]: # (This file is auto-generated. Please do not modify it yourself.)

# Protobuf Documentation

<a name="top"></a>

## Purpose

This module is designed to calculate an inflation rewards each years based on static params. Each block rewards is
the same over the year.

### 🧮 Calculation

The initial inflation is set to 15%, `annual_provisions` and `target_supply` will set at the beginning of the chain (on the first block).
It's based on the initial total supply. For example, with a total supply of 200M token, annual_provisions will be configured to 30M and target_supply 230M.

At the end of the year (the last block of the year, to be more precise), due to rounding imprecision, if the distributed
tokens for the last block of the year added with the actual tokens total supply is bigger than the `target_supply`,
only the difference to reach the `target_supply` will be minted. Conversely, in some case, due to rounding also,
the `target_supply` is not reached at the last block, but at the next block.

The new inflation, `annual_provisions` and `target_supply` is recalculated at the next block after the block that
reach the `target_supply`. To calculate the new inflation, get the current inflation multiplied by (`1 - annual_reduction_factor`),
then the new `annual_provisions` and `target_supply` is deducted based on the current total supply (that is the old `target_supply`)

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

Minter represents the minting state.

At the beginning of the chain (first block) the mint module will recalculate the `annual_provisions` and
`target_supply` based on the genesis total token supply and the inflation configured.
By default inflation is set to 15%. If the genesis total token supply is 200M token, the `annual_provision` will be 30M
and `target_supply` 230M.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `inflation` | [string](#string) |  | current annual inflation rate |
| `annual_provisions` | [string](#string) |  | current annual expected provisions |

<a name="mint.v1beta1.Params"></a>

### Params

Params holds parameters for the mint module.

Configure the annual reduction factor will update at the each end of year the new token distribution rate by reducing
the actual inflation by the `annual_reduction_factor` configured.
By default, `annual_reduction_factor` is 20%. For example, with an initial inflation of 15%, at the end of the year,
new inflation will be 12%.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `mint_denom` | [string](#string) |  | type of coin to mint |
| `inflation_coef` | [string](#string) |  | annual inflation coefficient |
| `bounding_adjustment` | [string](#string) |  | bounding adjustment |
| `target_bounding_ratio` | [string](#string) |  | target boudning ratio |
| `blocks_per_year` | [uint64](#uint64) |  | expected blocks per yearmake p |

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
