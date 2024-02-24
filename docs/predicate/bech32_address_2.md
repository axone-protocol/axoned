---
sidebar_position: 4
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# bech32_address/2

## Description

`bech32_address/2` is a predicate that convert a [bech32](<https://docs.cosmos.network/main/build/spec/addresses/bech32#hrp-table>) encoded string into [base64](<https://fr.wikipedia.org/wiki/Base64>) bytes and give the address prefix, or convert a prefix \(HRP\) and [base64](<https://fr.wikipedia.org/wiki/Base64>) encoded bytes to [bech32](<https://docs.cosmos.network/main/build/spec/addresses/bech32#hrp-table>) encoded string.

The signature is as follows:

```text
bech32_address(-Address, +Bech32)
bech32_address(+Address, -Bech32)
bech32_address(+Address, +Bech32)
```

where:

- Address is a pair of the HRP \(Human\-Readable Part\) which holds the address prefix and a list of numbers ranging from 0 to 255 that represent the base64 encoded bech32 address string.
- Bech32 is an Atom or string representing the bech32 encoded string address

## Examples

```text
# Convert the given bech32 address into base64 encoded byte by unify the prefix of given address (Hrp) and the
base64 encoded value (Address).
- bech32_address(-(Hrp, Address), 'okp415wn30a9z4uc692s0kkx5fp5d4qfr3ac7sj9dqn').

# Convert the given pair of HRP and base64 encoded address byte by unify the Bech32 string encoded value.
- bech32_address(-('okp4', [163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]), Bech32).
```
