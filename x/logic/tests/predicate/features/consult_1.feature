Feature: consult/1
  This feature is to test the consult/1 predicate.

  @great_for_documentation
  Scenario: Consult a Prolog program from the embedded library
  This scenario demonstrates how to load a library file and use one of its predicates.

    Given the query:
      """ prolog
      consult('/v1/lib/lists.pl'),
      member(Who, [alice,bob]).
      """
    When the query is run (limited to 1 solutions)
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4013
      answer:
        has_more: true
        variables: ["Who"]
        results:
        - substitutions:
          - variable: Who
            expression: "alice"
      """

  @great_for_documentation
  Scenario: Consult several Prolog programs at once
  This scenario demonstrates consult/1 with a list of files.

    Given the program:
      """ prolog
      :- consult([
        '/v1/lib/bank.pl',
        '/v1/lib/chain.pl'
      ]).
      """
    Given the query:
      """ prolog
      current_predicate(bank_balances/2),
      current_predicate(header_info/1).
      """
    When the query is run (limited to 2 solutions)
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4758
      answer:
        has_more: false
        variables:
        results:
        - substitutions:
      """
