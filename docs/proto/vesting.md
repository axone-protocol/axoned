[//]: # (This file is auto-generated. Please do not modify it yourself.)

# Protobuf Documentation

<a name="top"></a>

## Table of Contents

- [vesting/v1beta1/vesting.proto](#vesting/v1beta1/vesting.proto)
  - [BaseVestingAccount](#vesting.v1beta1.BaseVestingAccount)
  - [CliffVestingAccount](#vesting.v1beta1.CliffVestingAccount)
  - [ContinuousVestingAccount](#vesting.v1beta1.ContinuousVestingAccount)
  - [DelayedVestingAccount](#vesting.v1beta1.DelayedVestingAccount)
  - [Period](#vesting.v1beta1.Period)
  - [PeriodicVestingAccount](#vesting.v1beta1.PeriodicVestingAccount)
  - [PermanentLockedAccount](#vesting.v1beta1.PermanentLockedAccount)
  
- [vesting/v1beta1/tx.proto](#vesting/v1beta1/tx.proto)
  - [MsgCreateCliffVestingAccount](#vesting.v1beta1.MsgCreateCliffVestingAccount)
  - [MsgCreateCliffVestingAccountResponse](#vesting.v1beta1.MsgCreateCliffVestingAccountResponse)
  - [MsgCreatePeriodicVestingAccount](#vesting.v1beta1.MsgCreatePeriodicVestingAccount)
  - [MsgCreatePeriodicVestingAccountResponse](#vesting.v1beta1.MsgCreatePeriodicVestingAccountResponse)
  - [MsgCreatePermanentLockedAccount](#vesting.v1beta1.MsgCreatePermanentLockedAccount)
  - [MsgCreatePermanentLockedAccountResponse](#vesting.v1beta1.MsgCreatePermanentLockedAccountResponse)
  - [MsgCreateVestingAccount](#vesting.v1beta1.MsgCreateVestingAccount)
  - [MsgCreateVestingAccountResponse](#vesting.v1beta1.MsgCreateVestingAccountResponse)
  
  - [Msg](#vesting.v1beta1.Msg)
  
- [Scalar Value Types](#scalar-value-types)

<a name="vesting/v1beta1/vesting.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## vesting/v1beta1/vesting.proto

<a name="vesting.v1beta1.BaseVestingAccount"></a>

### BaseVestingAccount

BaseVestingAccount implements the VestingAccount interface. It contains all
the necessary fields needed for any vesting account implementation.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_account` | [cosmos.auth.v1beta1.BaseAccount](#cosmos.auth.v1beta1.BaseAccount) |  |  |
| `original_vesting` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `delegated_free` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `delegated_vesting` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `end_time` | [int64](#int64) |  | Vesting end time, as unix timestamp (in seconds). |

<a name="vesting.v1beta1.CliffVestingAccount"></a>

### CliffVestingAccount

CliffVestingAccount implements the VestingAccount interface. It
continuously vests by unlocking coins after a cliff period linearly with respect to time.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_vesting_account` | [BaseVestingAccount](#vesting.v1beta1.BaseVestingAccount) |  | base_vesting_account implements the VestingAccount interface. It contains all the necessary fields needed for any vesting account implementation |
| `start_time` | [int64](#int64) |  | start_time defines the time at which the vesting period begins |
| `cliff_time` | [int64](#int64) |  |  |

<a name="vesting.v1beta1.ContinuousVestingAccount"></a>

### ContinuousVestingAccount

ContinuousVestingAccount implements the VestingAccount interface. It
continuously vests by unlocking coins linearly with respect to time.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_vesting_account` | [BaseVestingAccount](#vesting.v1beta1.BaseVestingAccount) |  |  |
| `start_time` | [int64](#int64) |  | Vesting start time, as unix timestamp (in seconds). |

<a name="vesting.v1beta1.DelayedVestingAccount"></a>

### DelayedVestingAccount

DelayedVestingAccount implements the VestingAccount interface. It vests all
coins after a specific time, but non prior. In other words, it keeps them
locked until a specified time.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_vesting_account` | [BaseVestingAccount](#vesting.v1beta1.BaseVestingAccount) |  |  |

<a name="vesting.v1beta1.Period"></a>

### Period

Period defines a length of time and amount of coins that will vest.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `length` | [int64](#int64) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |

<a name="vesting.v1beta1.PeriodicVestingAccount"></a>

### PeriodicVestingAccount

PeriodicVestingAccount implements the VestingAccount interface. It
periodically vests by unlocking coins during each specified period.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_vesting_account` | [BaseVestingAccount](#vesting.v1beta1.BaseVestingAccount) |  |  |
| `start_time` | [int64](#int64) |  |  |
| `vesting_periods` | [Period](#vesting.v1beta1.Period) | repeated |  |

<a name="vesting.v1beta1.PermanentLockedAccount"></a>

### PermanentLockedAccount

PermanentLockedAccount implements the VestingAccount interface. It does
not ever release coins, locking them indefinitely. Coins in this account can
still be used for delegating and for governance votes even while locked.

Since: cosmos-sdk 0.43

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_vesting_account` | [BaseVestingAccount](#vesting.v1beta1.BaseVestingAccount) |  |  |

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

 [//]: # (end services)

<a name="vesting/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## vesting/v1beta1/tx.proto

<a name="vesting.v1beta1.MsgCreateCliffVestingAccount"></a>

### MsgCreateCliffVestingAccount

MsgCreateCliffVestingAccount defines a message that enables creating a cliff vesting
account.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from_address` | [string](#string) |  |  |
| `to_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `end_time` | [int64](#int64) |  |  |
| `cliff_time` | [int64](#int64) |  |  |

<a name="vesting.v1beta1.MsgCreateCliffVestingAccountResponse"></a>

### MsgCreateCliffVestingAccountResponse

MsgCreateCliffVestingAccountResponse defines the Msg/CreateVestingAccount response type.

<a name="vesting.v1beta1.MsgCreatePeriodicVestingAccount"></a>

### MsgCreatePeriodicVestingAccount

MsgCreateVestingAccount defines a message that enables creating a vesting
account.

Since: cosmos-sdk 0.46

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from_address` | [string](#string) |  |  |
| `to_address` | [string](#string) |  |  |
| `start_time` | [int64](#int64) |  | start of vesting as unix time (in seconds). |
| `vesting_periods` | [Period](#vesting.v1beta1.Period) | repeated |  |

<a name="vesting.v1beta1.MsgCreatePeriodicVestingAccountResponse"></a>

### MsgCreatePeriodicVestingAccountResponse

MsgCreateVestingAccountResponse defines the Msg/CreatePeriodicVestingAccount
response type.

Since: cosmos-sdk 0.46

<a name="vesting.v1beta1.MsgCreatePermanentLockedAccount"></a>

### MsgCreatePermanentLockedAccount

MsgCreatePermanentLockedAccount defines a message that enables creating a permanent
locked account.

Since: cosmos-sdk 0.46

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from_address` | [string](#string) |  |  |
| `to_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |

<a name="vesting.v1beta1.MsgCreatePermanentLockedAccountResponse"></a>

### MsgCreatePermanentLockedAccountResponse

MsgCreatePermanentLockedAccountResponse defines the Msg/CreatePermanentLockedAccount response type.

Since: cosmos-sdk 0.46

<a name="vesting.v1beta1.MsgCreateVestingAccount"></a>

### MsgCreateVestingAccount

MsgCreateVestingAccount defines a message that enables creating a vesting
account.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from_address` | [string](#string) |  |  |
| `to_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `end_time` | [int64](#int64) |  |  |
| `delayed` | [bool](#bool) |  |  |

<a name="vesting.v1beta1.MsgCreateVestingAccountResponse"></a>

### MsgCreateVestingAccountResponse

MsgCreateVestingAccountResponse defines the Msg/CreateVestingAccount response type.

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

<a name="vesting.v1beta1.Msg"></a>

### Msg

Msg defines the bank Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CreateVestingAccount` | [MsgCreateVestingAccount](#vesting.v1beta1.MsgCreateVestingAccount) | [MsgCreateVestingAccountResponse](#vesting.v1beta1.MsgCreateVestingAccountResponse) | CreateVestingAccount defines a method that enables creating a vesting account. | |
| `CreatePermanentLockedAccount` | [MsgCreatePermanentLockedAccount](#vesting.v1beta1.MsgCreatePermanentLockedAccount) | [MsgCreatePermanentLockedAccountResponse](#vesting.v1beta1.MsgCreatePermanentLockedAccountResponse) | CreatePermanentLockedAccount defines a method that enables creating a permanent locked account.

Since: cosmos-sdk 0.46 | |
| `CreatePeriodicVestingAccount` | [MsgCreatePeriodicVestingAccount](#vesting.v1beta1.MsgCreatePeriodicVestingAccount) | [MsgCreatePeriodicVestingAccountResponse](#vesting.v1beta1.MsgCreatePeriodicVestingAccountResponse) | CreatePeriodicVestingAccount defines a method that enables creating a periodic vesting account.

Since: cosmos-sdk 0.46 | |
| `CreateCliffVestingAccount` | [MsgCreateCliffVestingAccount](#vesting.v1beta1.MsgCreateCliffVestingAccount) | [MsgCreateCliffVestingAccountResponse](#vesting.v1beta1.MsgCreateCliffVestingAccountResponse) | CreateCliffVestingAccount defines a method that enables creating a cliff vesting account. | |

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
