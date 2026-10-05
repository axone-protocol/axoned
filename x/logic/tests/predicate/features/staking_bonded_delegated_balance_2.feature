Feature: staking_bonded_delegated_balance/2
  Sum only delegations to bonded validators with exact coin arithmetic.

  Background:
    Given the program:
      """ prolog
      :- consult('/v1/lib/staking.pl').
      """
    And the staking bond denomination is "ustake"
    And the staking validators:
      | operator                | status    | jailed |
      | axonevaloper1bonded      | bonded    | false  |
      | axonevaloper1jailed      | bonded    | true   |
      | axonevaloper1unbonding   | unbonding | false  |
      | axonevaloper1unbonded    | unbonded  | false  |

  @great_for_documentation
  Scenario: Sum bonded delegations including a jailed bonded validator
    Given the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following staking delegations:
      | validator               | denom  | amount |
      | axonevaloper1unbonding   | ustake | 900    |
      | axonevaloper1jailed      | ustake | 75     |
      | axonevaloper1unbonded    | ustake | 800    |
      | axonevaloper1bonded      | ustake | 25     |
    And the query:
      """ prolog
      staking_bonded_delegated_balance('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', Balance).
      """
    When the query is run
    Then the logical answer we get is:
      """ yaml
      answer:
        has_more: false
        variables: [Balance]
        results:
        - substitutions:
          - variable: Balance
            expression: "coin(ustake,100)"
      """

  Scenario Outline: Preserve integer boundaries and decimal precision beyond int64
    Given the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following staking delegations:
      | validator             | denom  | amount  |
      | axonevaloper1bonded    | ustake | <Left>  |
      | axonevaloper1jailed    | ustake | <Right> |
    And the query:
      """ prolog
      staking_bonded_delegated_balance('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', Balance).
      """
    When the query is run
    Then the logical answer we get is:
      """ yaml
      answer:
        has_more: false
        variables: [Balance]
        results:
        - substitutions:
          - variable: Balance
            expression: "coin(ustake,<Total>)"
      """
    Examples:
      | Left                 | Right               | Total                  |
      | 9223372036854775806  | 1                   | 9223372036854775807     |
      | 9223372036854775807  | 1                   | '9223372036854775808'   |
      | 9223372036854775808  | 9223372036854775808  | '18446744073709551616'  |
      | 99999999999999999999 | 1                   | '100000000000000000000' |

  Scenario: Return zero with the staking bond denomination when none are eligible
    Given the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following staking delegations:
      | validator               | denom  | amount |
      | axonevaloper1unbonding   | ustake | 900    |
    And the query:
      """ prolog
      staking_bonded_delegated_balance('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', Balance).
      """
    When the query is run
    Then the logical answer we get is:
      """ yaml
      answer:
        has_more: false
        variables: [Balance]
        results:
        - substitutions:
          - variable: Balance
            expression: "coin(ustake,0)"
      """

  Scenario: Return zero with the staking bond denomination for no delegations
    Given the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following staking delegations:
      | validator | denom | amount |
    And the query:
      """ prolog
      staking_bonded_delegated_balance('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', Balance).
      """
    When the query is run
    Then the logical answer we get is:
      """ yaml
      answer:
        has_more: false
        variables: [Balance]
        results:
        - substitutions:
          - variable: Balance
            expression: "coin(ustake,0)"
      """

  Scenario: Reject an eligible coin with a different denomination
    Given the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following staking delegations:
      | validator             | denom  | amount |
      | axonevaloper1bonded    | uother | 50     |
    And the query:
      """ prolog
      staking_bonded_delegated_balance('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', _).
      """
    When the query is run
    Then the logical answer we get is:
      """ yaml
      answer:
        has_more: false
        results:
        - error: "error(domain_error(staking_bond_denom,uother),staking_bonded_delegated_balance/2)"
      """
