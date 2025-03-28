---
sidebar_position: 10
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# bech32_address/2

## Description

`bech32_address/2` is a predicate that converts a Bech32\-encoded string into a prefix \(HRP\) and Base64\-encoded bytes, or constructs a Bech32\-encoded string from a prefix and Base64 bytes.

This predicate handles Bech32 address encoding and decoding as per the Cosmos specification. In the Cosmos ecosystem, most chains \(e.g., Cosmos Hub, Akash\) share the BIP\-44 coin type 118', allowing HRP conversion \(e.g., 'cosmos' to 'akash'\) to produce valid addresses from the same underlying key.

## Signature

```text
bech32_address(-Address, +Bech32) is det
bech32_address(+Address, -Bech32) is det
```

where:

- Address: A pair \`HRP\-Base64Bytes\`, where: HRP is an atom representing the Human\-Readable Part \(e.g. 'cosmos', 'akash', 'axone'\), and Base64Bytes is a list of integers \(0\-255\) representing the Base64\-encoded bytes git statof the address.
- Bech32: An atom or string representing the Bech32\-encoded address \(e.g., 'cosmos17sc02mcgjzdv5l4jwnzffxw7g60y5ta4pggcp4'\).

## Limitations

Conversion between HRPs is only valid for chains sharing the same BIP\-44 coin type \(e.g., 118'\). For chains with distinct coin types \(e.g., Secret: 529', Bitsong: 639'\), this predicate cannot derive the correct address from another chain’s Bech32 string.

## References

- [Bech32 on Cosmos](<https://docs.cosmos.network/main/build/spec/addresses/bech32>)

- [Base64 Encoding](<https://fr.wikipedia.org/wiki/Base64>)

- [Cosmos Chain Registry](<https://github.com/cosmos/chain-registry>)

- [BIP 44](<https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki>)

## Examples

### Decode Bech32 Address into its Address Pair representation

This scenario demonstrates how to parse a provided bech32 address string into its `Address` pair representation.
An `Address` is a compound term `-` with two arguments, the first being the human-readable part (Hrp) and the second
being the numeric address as a list of integers ranging from 0 to 255 representing the bytes of the address in
base 64.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
bech32_address(Address, 'axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3975
answer:
  has_more: false
  variables: ["Address"]
  results:
  - substitutions:
    - variable: Address
      expression: "axone-[163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]"
```

### Decode Hrp and Address from a bech32 address

This scenario illustrates how to decode a bech32 address into the human-readable part (Hrp) and the numeric address.
The process extracts these components from a given bech32 address string, showcasing the ability to parse and
separate the address into its constituent parts.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
bech32_address(-(Hrp, Address), 'axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3975
answer:
  has_more: false
  variables: ["Hrp", "Address"]
  results:
  - substitutions:
    - variable: Hrp
      expression: "axone"
    - variable: Address
      expression: "[163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]"
```

### Extract Address only for AXONE bech32 address

This scenario demonstrates how to extract the address from a bech32 address string, specifically for a known
protocol, in this case, the axone protocol.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
bech32_address(-(axone, Address), 'axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3975
answer:
  has_more: false
  variables: ["Address"]
  results:
  - substitutions:
    - variable: Address
      expression: "[163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]"
```

### Encode Address Pair into Bech32 Address

This scenario demonstrates how to encode an `Address` pair representation into a bech32 address string.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
bech32_address(-('axone', [163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]), Bech32).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3975
answer:
  has_more: false
  variables: ["Bech32"]
  results:
  - substitutions:
    - variable: Bech32
      expression: "axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4"
```

### Check if a bech32 address is part of the axone protocol

This scenario shows how to check if a bech32 address is part of the axone protocol.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
axone_addr(Addr) :- bech32_address(-('axone', _), Addr).
```

- **Given** the query:

```  prolog
axone_addr('axone1p8u47en82gmzfm259y6z93r9qe63l25d858vqu').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3976
answer:
  has_more: false
  results:
  - substitutions:
```

### Error on Incorrect Bech32 Address format

This scenario demonstrates the system's response to an incorrect bech32 address format.
In this case, the system generates a `domain_error`, indicating that the provided argument does not meet the
expected format for a bech32 address.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
bech32_address(Address, axoneincorrect).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3975
answer:
  has_more: false
  variables: ["Address"]
  results:
  - error: "error(domain_error(encoding(bech32),axoneincorrect),[d,e,c,o,d,i,n,g, ,b,e,c,h,3,2, ,f,a,i,l,e,d,:, ,i,n,v,a,l,i,d, ,s,e,p,a,r,a,t,o,r, ,i,n,d,e,x, ,-,1],bech32_address/2)"
```

### Error on Incorrect Bech32 Address type

This scenario demonstrates the system's response to an incorrect bech32 address type.
In this case, the system generates a `type_error`, indicating that the provided argument does not meet the
expected type.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
bech32_address(-('axone', X), foo(bar)).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3975
answer:
  has_more: false
  variables: ["X"]
  results:
  - error: "error(type_error(atom,foo(bar)),bech32_address/2)"
```
