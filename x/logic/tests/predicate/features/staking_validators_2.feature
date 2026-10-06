Feature: staking_validators/2
  Query validators by staking bond status, independently of jailed status.

  Background:
    Given the program:
      """ prolog
      :- consult('/v1/lib/staking.pl').
      """
    And the staking validators:
      | operator                | status    | jailed |
      | axonevaloper1bonded      | bonded    | false  |
      | axonevaloper1jailed      | bonded    | true   |
      | axonevaloper1unbonding   | unbonding | false  |
      | axonevaloper1unbonded    | unbonded  | false  |

  @great_for_documentation
  Scenario: Select bonded validators including jailed validators
    Given the query:
      """ prolog
      staking_validators(bonded, Validators).
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

  Scenario Outline: Select non-bonded validators with the requested bond status
    Given the query:
      """ prolog
      staking_validators(<Status>, Validators).
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
            expression: "<Validators>"
      """
    Examples:
      | Status    | Validators                                                                                                        |
      | unbonding | [validator{operator:axonevaloper1unbonding,status:unbonding}]                                                       |
      | unbonded  | [validator{operator:axonevaloper1unbonded,status:unbonded}]                                                         |

  Scenario Outline: Reject a missing or invalid bond status
    Given the query:
      """ prolog
      staking_validators(<Status>, _).
      """
    When the query is run
    Then the logical answer we get is:
      """ yaml
      answer:
        has_more: false
        results:
        - error: "<Error>"
      """
    Examples:
      | Status   | Error                                                                                 |
      | active   | error(type_error(oneof([bonded,unbonding,unbonded]),active),staking_validators/2)        |
      | 42       | error(type_error(atom,42),staking_validators/2)                                         |
      | _        | error(instantiation_error,staking_validators/2)                                        |
