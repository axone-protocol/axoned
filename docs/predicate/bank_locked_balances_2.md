---
sidebar_position: 2
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# bank_locked_balances/2

## Description

`bank_locked_balances/2` is a predicate which unifies the given terms with the list of locked coins of the given account.

The signature is as follows:

```text
bank_locked_balances(?Account, ?Balances)
```

where:

- Account represents the account address \(in Bech32 format\).
- Balances represents the locked balances of the account as a list of pairs of coin denomination and amount.

## Examples

```text
# Query the locked coins of the account.
- bank_locked_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', X).

# Query the locked balances of all accounts. The result is a list of pairs of account address and balances.
- bank_locked_balances(X, Y).

# Query the first locked balances of the given account by unifying the denomination and amount with the given terms.
- bank_locked_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', [-(D, A), _]).
```
