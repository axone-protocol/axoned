Feature: staking_delegations_to_bonded_validators/2
  Filter delegations by validator bond status without filtering jailed validators.

  Background:
    Given the program:
      """ prolog
      :- consult('/v1/lib/staking.pl').
      """
    And the staking validators:
      | operator                | status    | jailed |
      | axonevaloper1bonded      | bonded    | false  |
      | axonevaloper1jailed      | bonded    | true   |
      | axonevaloper1unused      | bonded    | false  |
      | axonevaloper1unbonding   | unbonding | false  |
      | axonevaloper1unbonded    | unbonded  | false  |

  @great_for_documentation
  Scenario: Keep balances and order for bonded validators only
    Given the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following staking delegations:
      | validator               | denom  | amount |
      | axonevaloper1unbonding   | uaxone | 900    |
      | axonevaloper1jailed      | uaxone | 75     |
      | axonevaloper1unbonded    | uaxone | 800    |
      | axonevaloper1bonded      | uaxone | 25     |
      | axonevaloper1absent      | uaxone | 700    |
    And the query:
      """ prolog
      staking_delegations_to_bonded_validators('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', Delegations).
      """
    When the query is run
    Then the logical answer we get is:
      """ yaml
      answer:
        has_more: false
        variables: [Delegations]
        results:
        - substitutions:
          - variable: Delegations
            expression: "[delegation{balance:coin(uaxone,75),validator:axonevaloper1jailed},delegation{balance:coin(uaxone,25),validator:axonevaloper1bonded}]"
      """

  Scenario: Return an empty list for an account without delegations
    Given the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following staking delegations:
      | validator | denom | amount |
    And the query:
      """ prolog
      staking_delegations_to_bonded_validators('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', Delegations).
      """
    When the query is run
    Then the logical answer we get is:
      """ yaml
      answer:
        has_more: false
        variables: [Delegations]
        results:
        - substitutions:
          - variable: Delegations
            expression: "[]"
      """

  Scenario: Reject a delegator that is not an atom
    Given the query:
      """ prolog
      staking_delegations_to_bonded_validators(42, _).
      """
    When the query is run
    Then the logical answer we get is:
      """ yaml
      answer:
        has_more: false
        results:
        - error: "error(type_error(atom,42),staking_delegations_to_bonded_validators/2)"
      """
