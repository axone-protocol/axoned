Feature: staking_delegations/2
  Query all staking delegations of an account, regardless of validator bond status.

  @great_for_documentation
  Scenario: Query all delegations of an account
    Given the program:
      """ prolog
      :- consult('/v1/lib/staking.pl').
      """
    And the account "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa" has the following staking delegations:
      | validator             | denom  | amount  |
      | axonevaloper1validator | uaxone | 1250000 |
    And the staking validators:
      | operator              | status   | jailed |
      | axonevaloper1validator | unbonded | false  |
    Given the query:
      """ prolog
      staking_delegations('axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa', Delegations).
      """
    When the query is run
    Then the logical answer we get is:
      """ yaml
      answer:
        has_more: false
        variables: ["Delegations"]
        results:
        - substitutions:
          - variable: Delegations
            expression: "[delegation{balance:coin(uaxone,1250000),validator:axonevaloper1validator}]"
      """

  Scenario: Fail when the delegator is not an atom
    Given the program:
      """ prolog
      :- consult('/v1/lib/staking.pl').
      """
    Given the query:
      """ prolog
      staking_delegations(42, _).
      """
    When the query is run
    Then the logical answer we get is:
      """ yaml
      answer:
        has_more: false
        results:
        - error: "error(type_error(atom,42),staking_delegations/2)"
      """
