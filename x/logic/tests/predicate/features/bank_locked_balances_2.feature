Feature: bank_locked_balances/2
  This feature is to test the bank_locked_balances/2 predicate.

  @great_for_documentation
  Scenario: Query locked balances of an account with coins
    This scenario demonstrates how to query the locked balances of an account.

    Given the program:
      """ prolog
      :- consult('/v1/lib/bank.pl').
      """
    And the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following locked balances:
      | denom  | amount |
      | uaxone | 300    |
      | uatom  | 125    |
    Given the query:
      """ prolog
      bank_locked_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', Balances).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4033
      answer:
        has_more: false
        variables: ["Balances"]
        results:
        - substitutions:
          - variable: Balances
            expression: "[uatom-125,uaxone-300]"
      """

  @great_for_documentation
  Scenario: Query locked balances of an account with no coins
    This scenario demonstrates querying locked balances for an account that currently has no locked coin.

    Given the program:
      """ prolog
      :- consult('/v1/lib/bank.pl').
      """
    And the account "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep" has the following locked balances:
      | denom | amount |
    Given the query:
      """ prolog
      bank_locked_balances('axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep', Balances).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4017
      answer:
        has_more: false
        variables: ["Balances"]
        results:
        - substitutions:
          - variable: Balances
            expression: "[]"
      """

  @great_for_documentation
  Scenario: Use member/2 to check for a specific locked coin
    This scenario demonstrates using member/2 to retrieve one specific locked denomination.

    Given the program:
      """ prolog
      :- consult('/v1/lib/bank.pl').

      locked_has_coin(Address, Denom, Amount) :-
          bank_locked_balances(Address, Balances),
          member(Denom-Amount, Balances).
      """
    And the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following locked balances:
      | denom  | amount |
      | uaxone | 1000   |
      | uatom  | 500    |
    Given the query:
      """ prolog
      locked_has_coin('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', uaxone, Amount).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4037
      answer:
        has_more: false
        variables: ["Amount"]
        results:
        - substitutions:
          - variable: Amount
            expression: "1000"
      """

  Scenario: Query locked balances with amount greater than int64
    This scenario validates that very large locked amounts are preserved as atoms.

    Given the program:
      """ prolog
      :- consult('/v1/lib/bank.pl').
      """
    And the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following locked balances:
      | denom  | amount              |
      | uaxone | 9223372036854775808 |
    Given the query:
      """ prolog
      bank_locked_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', [uaxone-'9223372036854775808']).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4025
      answer:
        has_more: false
        results:
        - {}
      """

  @great_for_documentation
  Scenario: Fail when address is a variable
    This scenario shows what happens when the address argument is left unbound.

    Given the program:
      """ prolog
      :- consult('/v1/lib/bank.pl').
      """
    Given the query:
      """ prolog
      bank_locked_balances(Address, Balances).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4002
      answer:
        has_more: false
        variables: ["Address", "Balances"]
        results:
        - error: "error(instantiation_error,must_be/2)"
      """

  @great_for_documentation
  Scenario: Fail with invalid address format
    This scenario shows the error returned when the address is not a valid Bech32 value.

    Given the program:
      """ prolog
      :- consult('/v1/lib/bank.pl').
      """
    Given the query:
      """ prolog
      bank_locked_balances('invalid_address', Balances).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3999
      answer:
        has_more: false
        variables: ["Balances"]
        results:
        - error: "error(domain_error(encoding(bech32),invalid_address),bank_locked_balances/2)"
      """

  @great_for_documentation
  Scenario: Fail when address is not an atom
    This scenario shows that the address must be an atom (e.g. a Bech32 string), not a number.

    Given the program:
      """ prolog
      :- consult('/v1/lib/bank.pl').
      """
    Given the query:
      """ prolog
      bank_locked_balances(42, _).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3999
      answer:
        has_more: false
        results:
        - error: "error(type_error(atom,42),bech32_address/2)"
      """
