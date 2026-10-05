Feature: staking_bonded_validators/1
  Return only validators whose staking status is bonded.

  Background:
    Given the program:
      """ prolog
      :- consult('/v1/lib/staking.pl').
      """

  Scenario: Include jailed validators if their bond status is still bonded
    Given the staking validators:
      | operator                | status    | jailed |
      | axonevaloper1bonded      | bonded    | false  |
      | axonevaloper1jailed      | bonded    | true   |
      | axonevaloper1unbonding   | unbonding | false  |
    And the query:
      """ prolog
      staking_bonded_validators(Validators).
      """
    When the query is run
    Then the logical answer we get is:
      """ yaml
      answer:
        has_more: false
        variables: [Validators]
        results:
        - substitutions:
          - variable: Validators
            expression: "[validator{operator:axonevaloper1bonded,status:bonded},validator{operator:axonevaloper1jailed,status:bonded}]"
      """

  Scenario: Return an empty list when no validator is bonded
    Given the staking validators:
      | operator                | status    | jailed |
      | axonevaloper1unbonded    | unbonded  | false  |
    And the query:
      """ prolog
      staking_bonded_validators(Validators).
      """
    When the query is run
    Then the logical answer we get is:
      """ yaml
      answer:
        has_more: false
        variables: [Validators]
        results:
        - substitutions:
          - variable: Validators
            expression: "[]"
      """
