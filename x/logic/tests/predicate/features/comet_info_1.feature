Feature: comet_info/1
  This feature is to test the comet_info/1 predicate.

  @great_for_documentation
  Scenario: Retrieve current comet block info
    This scenario demonstrates how to retrieve the current Comet block information available to the query.
    The Comet info contains consensus-related metadata such as proposer address, validators hash,
    evidence, and last commit information.

    Given the query:
      """ prolog
      consult('/v1/lib/chain.pl'),
      comet_info(CometInfo).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3980
      answer:
        has_more: false
        variables: ["CometInfo"]
        results:
        - substitutions:
          - variable: CometInfo
            expression: >-
              comet{evidence:[],last_commit:commit_info{round:0,votes:[]},proposer_address:[],validators_hash:[]}
      """

  @great_for_documentation
  Scenario: Retrieve proposer address from current comet block info
    This scenario demonstrates how to read a specific field from comet_info/1.

    Given the program:
      """ prolog
      :- consult('/v1/lib/chain.pl').

      proposer_address(Address) :-
          comet_info(CometInfo),
          Address = CometInfo.proposer_address.
      """
    Given the query:
      """ prolog
      proposer_address(Address).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3983
      answer:
        has_more: false
        variables: ["Address"]
        results:
        - substitutions:
          - variable: Address
            expression: "[]"
      """

  @great_for_documentation
  Scenario: Retrieve last commit round from current comet block info
    This scenario demonstrates how to read nested fields from comet_info/1.

    Given the program:
      """ prolog
      :- consult('/v1/lib/chain.pl').

      last_commit_round(Round) :-
          comet_info(CometInfo),
          Round = CometInfo.last_commit.round.
      """
    Given the query:
      """ prolog
      last_commit_round(Round).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3984
      answer:
        has_more: false
        variables: ["Round"]
        results:
        - substitutions:
          - variable: Round
            expression: "0"
      """

  @great_for_documentation
  Scenario: Evaluate a condition based on comet evidence and commit round
    This scenario demonstrates how to combine comet_info/1 fields in a rule.
    It checks that there is no evidence and that the last commit round is non-negative.

    Given the program:
      """ prolog
      :- consult('/v1/lib/chain.pl').

      healthy_consensus_snapshot :-
          comet_info(CometInfo),
          CometInfo.evidence = [],
          CometInfo.last_commit.round >= 0,
          !.
      """
    Given the query:
      """ prolog
      healthy_consensus_snapshot.
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3986
      answer:
        has_more: false
        results:
        - substitutions:
      """
