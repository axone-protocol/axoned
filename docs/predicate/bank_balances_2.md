---
sidebar_position: 15
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# bank_balances/2

## Module

This predicate is provided by `bank.pl`.

Load this module before using the predicate:

```prolog
:- consult('/v1/lib/bank.pl').
```

## Description

Unifies Balances with the list of coin balances for the given account Address.
The address must be instantiated (non-variable) and in Bech32 format.

Returned term shape:

```prolog
[Denom-Amount, ...]
```

where:

- Denom is an atom representing the coin denomination.
- Amount is an integer when it fits in int64, otherwise an atom preserving full precision.
- The list is sorted by denomination.

Throws instantiation_error if Address is a variable.
Throws domain_error(encoding(bech32), Address) if Address is not a valid Bech32 address.

Examples:

```prolog
?- bank_balances('axone1...', Balances).
Balances = [uatom-100, uaxone-200].
```

## Signature

```text
bank_balances(+Address, -Balances) is det
```

## Examples

### Query balances of an account with coins

This scenario demonstrates how to query the balances of an account that holds multiple coins.
The bank_balances/2 predicate retrieves all coin balances for a given account address.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- consult('/v1/lib/bank.pl').
```

- **And** the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following balances:

| key | value |
| --- | ----- |
| denom | amount |
| uaxone | 1000 |
| uatom | 500 |

- **Given** the query:

```  prolog
bank_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', Balances).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4033
answer:
  has_more: false
  variables: ["Balances"]
  results:
  - substitutions:
    - variable: Balances
      expression: "[uatom-500,uaxone-1000]"
```

### Query balances of an account with no coins

This scenario demonstrates querying the balances of an account that exists but has no coins.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- consult('/v1/lib/bank.pl').
```

- **And** the account "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep" has the following balances:

| key | value |
| --- | ----- |
| denom | amount |

- **Given** the query:

```  prolog
bank_balances('axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep', Balances).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4017
answer:
  has_more: false
  variables: ["Balances"]
  results:
  - substitutions:
    - variable: Balances
      expression: "[]"
```

### Query balances and extract first coin denomination

This scenario demonstrates pattern matching on the balances list to extract specific coin information.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- consult('/v1/lib/bank.pl').

first_denom(Address, Denom) :-
    bank_balances(Address, Balances),
    Balances = [Denom-_ | _].
```

- **And** the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following balances:

| key | value |
| --- | ----- |
| denom | amount |
| uaxone | 1000 |
| uatom | 500 |

- **Given** the query:

```  prolog
first_denom('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', Denom).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4035
answer:
  has_more: false
  variables: ["Denom"]
  results:
  - substitutions:
    - variable: Denom
      expression: "uatom"
```

### Use member/2 to check for a specific coin

This scenario demonstrates using member/2 to check if an account has a specific coin denomination.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- consult('/v1/lib/bank.pl').

has_coin(Address, Denom, Amount) :-
    bank_balances(Address, Balances),
    member(Denom-Amount, Balances).
```

- **And** the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following balances:

| key | value |
| --- | ----- |
| denom | amount |
| uaxone | 1000 |
| uatom | 500 |

- **Given** the query:

```  prolog
has_coin('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', uaxone, Amount).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4037
answer:
  has_more: false
  variables: ["Amount"]
  results:
  - substitutions:
    - variable: Amount
      expression: "1000"
```

### Fail when address is a variable

This scenario shows what happens when the address argument is left unbound.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- consult('/v1/lib/bank.pl').
```

- **Given** the query:

```  prolog
bank_balances(Address, Balances).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4002
answer:
  has_more: false
  variables: ["Address", "Balances"]
  results:
  - error: "error(instantiation_error,must_be/2)"
```

### Fail with invalid address format

This scenario shows the error returned when the address is not a valid Bech32 value.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- consult('/v1/lib/bank.pl').
```

- **Given** the query:

```  prolog
bank_balances('invalid_address', Balances).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3999
answer:
  has_more: false
  variables: ["Balances"]
  results:
  - error: "error(domain_error(encoding(bech32),invalid_address),bank_balances/2)"
```

### Fail when address is not an atom

This scenario shows that the address must be an atom (e.g. a Bech32 string), not a number.

Here are the steps of the scenario:

- **Given** the program:

```  prolog
:- consult('/v1/lib/bank.pl').
```

- **Given** the query:

```  prolog
bank_balances(42, _).
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 3999
answer:
  has_more: false
  results:
  - error: "error(type_error(atom,42),bech32_address/2)"
```
