Feature: bank_balances/2
  This feature is to test the bank_balances/2 predicate.

  @great_for_documentation
  Scenario: Query balances of an account with coins
    This scenario demonstrates how to query the balances of an account that holds multiple coins.
    The bank_balances/2 predicate retrieves all coin balances for a given account address.

    Given the program:
      """ prolog
      :- consult('/v1/lib/bank.pl').
      """
    And the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following balances:
      | denom  | amount |
      | uaxone | 1000   |
      | uatom  | 500    |
    Given the query:
      """ prolog
      bank_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', Balances).
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
            expression: "[uatom-500,uaxone-1000]"
      """

  @great_for_documentation
  Scenario: Query balances of an account with no coins
    This scenario demonstrates querying the balances of an account that exists but has no coins.

    Given the program:
      """ prolog
      :- consult('/v1/lib/bank.pl').
      """
    And the account "axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep" has the following balances:
      | denom | amount |
    Given the query:
      """ prolog
      bank_balances('axone1wze8mn5nsgl9qrgazq6a92fvh7m5e6ps372aep', Balances).
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
  Scenario: Query balances and extract first coin denomination
    This scenario demonstrates pattern matching on the balances list to extract specific coin information.

    Given the program:
      """ prolog
      :- consult('/v1/lib/bank.pl').

      first_denom(Address, Denom) :-
          bank_balances(Address, Balances),
          Balances = [Denom-_ | _].
      """
    And the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following balances:
      | denom  | amount |
      | uaxone | 1000   |
      | uatom  | 500    |
    Given the query:
      """ prolog
      first_denom('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', Denom).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4035
      answer:
        has_more: false
        variables: ["Denom"]
        results:
        - substitutions:
          - variable: Denom
            expression: "uatom"
      """

  @great_for_documentation
  Scenario: Use member/2 to check for a specific coin
    This scenario demonstrates using member/2 to check if an account has a specific coin denomination.

    Given the program:
      """ prolog
      :- consult('/v1/lib/bank.pl').

      has_coin(Address, Denom, Amount) :-
          bank_balances(Address, Balances),
          member(Denom-Amount, Balances).
      """
    And the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following balances:
      | denom  | amount |
      | uaxone | 1000   |
      | uatom  | 500    |
    Given the query:
      """ prolog
      has_coin('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', uaxone, Amount).
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

  Scenario: Query balances with IBC denomination
    This scenario validates that denominations requiring quoting (e.g. ibc/<hash>) round-trip correctly.

    Given the program:
      """ prolog
      :- consult('/v1/lib/bank.pl').

      has_ibc_coin(Address, Amount) :-
          bank_balances(Address, Balances),
          member('ibc/0123456789ABCDEF'-Amount, Balances).
      """
    And the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following balances:
      | denom               | amount |
      | ibc/0123456789ABCDEF | 777    |
      | uaxone              | 1000   |
    Given the query:
      """ prolog
      has_ibc_coin('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', Amount).
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
            expression: "777"
      """

  Scenario: Query balances with amount greater than int64
    This scenario validates that very large amounts are preserved as atoms instead of causing parser errors.

    Given the program:
      """ prolog
      :- consult('/v1/lib/bank.pl').
      """
    And the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following balances:
      | denom  | amount              |
      | uaxone | 9223372036854775808 |
    Given the query:
      """ prolog
      bank_balances('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', [uaxone-'9223372036854775808']).
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
      bank_balances(Address, Balances).
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
      bank_balances('invalid_address', Balances).
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
        - error: "error(domain_error(encoding(bech32),invalid_address),bank_balances/2)"
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
      bank_balances(42, _).
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
