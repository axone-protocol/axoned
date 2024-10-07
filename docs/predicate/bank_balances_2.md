---
sidebar_position: 1
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# bank_balances/2

## Description

`bank_balances/2` is a predicate which unifies the given terms with the list of balances \(coins\) of the given account.

The signature is as follows:

```text
bank_balances(?Address, ?Balances)
```

where:

- Address represents the account address \(in Bech32 format\).
- Balances represents the balances of the account as a list of pairs of coin denomination and amount.

## Examples

```text
# Query the balances of the account.
- bank_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', X).

# Query the balances of all accounts. The result is a list of pairs of account address and balances.
- bank_balances(X, Y).

# Query the first balance of the given account by unifying the denomination and amount with the given terms.
- bank_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', [-(D, A), _]).
```
