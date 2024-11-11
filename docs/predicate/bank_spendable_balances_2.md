---
sidebar_position: 8
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# bank_spendable_balances/2

## Description

`bank_spendable_balances/2` is a predicate which unifies the given terms with the list of spendable coins of the given account.

The signature is as follows:

```text
bank_spendable_balances(?Address, ?Balances)
```

where:

- Address represents the account address \(in Bech32 format\).
- Balances represents the spendable balances of the account as a list of pairs of coin denomination and amount.

## Examples

```text
# Query the spendable balances of the account.
- bank_spendable_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', X).

# Query the spendable balances of all accounts. The result is a list of pairs of account address and balances.
- bank_spendable_balances(X, Y).

# Query the first spendable balances of the given account by unifying the denomination and amount with the given terms.
- bank_spendable_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', [-(D, A), _]).
```
