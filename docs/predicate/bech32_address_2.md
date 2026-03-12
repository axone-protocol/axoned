---
sidebar_position: 21
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# bech32_address/2

## Module

This predicate is provided by `bech32.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/bech32.pl').
```

## Description

Converts between a Bech32 atom and its Address pair representation.

The predicate follows a functional direction:

- when Address is ground, it encodes Address into Bech32;
- otherwise, when Bech32 is ground, it decodes Bech32 into Address;
- otherwise, it throws instantiation_error.

Address is represented as Hrp-Bytes where:

- Hrp is an atom
- Bytes is a proper list of byte integers in [0,255]

## Signature

```text
bech32_address(?Address, ?Bech32) is det
```

## Examples

### Decode Bech32 Address into its Address Pair representation

This scenario demonstrates how to parse a provided bech32 address string into its `Address` pair representation.
An `Address` is a compound term `-` with two arguments, the first being the human-readable part (Hrp) and the second
being the numeric address as a list of integers ranging from 0 to 255 representing the bytes of the address in
base 64.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/bech32.pl'),
bech32_address(Address, 'axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 8474
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
consult('/v1/lib/bech32.pl'),
bech32_address(-(Hrp, Address), 'axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 8496
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
consult('/v1/lib/bech32.pl'),
bech32_address(-(axone, Address), 'axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 8496
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
consult('/v1/lib/bech32.pl'),
bech32_address(-('axone', [163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]), Bech32).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 33897
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
consult('/v1/lib/bech32.pl'),
axone_addr('axone1p8u47en82gmzfm259y6z93r9qe63l25d858vqu').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 8509
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
consult('/v1/lib/bech32.pl'),
bech32_address(Address, axoneincorrect).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 8797
answer:
  has_more: false
  variables: ["Address"]
  results:
  - error: "error(domain_error(valid_encoding(bech32),axoneincorrect),bech32_address/2)"
```

### Error on Incorrect Bech32 Address type

This scenario demonstrates the system's response to an incorrect bech32 address type.
In this case, the system generates a `type_error`, indicating that the provided argument does not meet the
expected type.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/bech32.pl'),
bech32_address(-('axone', X), foo(bar)).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4527
answer:
  has_more: false
  variables: ["X"]
  results:
  - error: "error(type_error(atom,foo(bar)),bech32_address/2)"
```
