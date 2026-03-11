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
      gas_used: 4792
      answer:
        has_more: false
        variables:
        results:
        - substitutions:
      """

  @great_for_documentation
  Scenario: Consult a published user Prolog library from the logic virtual file system
  This scenario demonstrates how to load a user Prolog library through the user-scoped publication view.

    Given the user Prolog library published by "axone15mefcxeleeefp2ga8yrax9tdzw7jkecjxeg7st" is:
      """ prolog
      member_lib(alice).
      """
    Given the query:
      """ prolog
      consult('/v1/var/lib/logic/users/axone15mefcxeleeefp2ga8yrax9tdzw7jkecjxeg7st/programs/42f889e07ab07b4764f19207799046cb603b954659b601d1a1238aaeac111d5d.pl'),
      member_lib(Who).
      """
    When the query is run (limited to 1 solutions)
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 14072
      answer:
        has_more: false
        variables: ["Who"]
        results:
        - substitutions:
          - variable: Who
            expression: "alice"
      """

  @great_for_documentation
  Scenario: Consult a non published user Prolog library
  This scenario demonstrates the error returned when the user-scoped publication view does not point to a published user Prolog library.

    Given the query:
      """ prolog
      consult('/v1/var/lib/logic/users/axone15mefcxeleeefp2ga8yrax9tdzw7jkecjxeg7st/programs/42f889e07ab07b4764f19207799046cb603b954659b601d1a1238aaeac111d5d.pl').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5219
      answer:
        has_more: false
        results:
        - error: "error(existence_error(source_sink,/v1/var/lib/logic/users/axone15mefcxeleeefp2ga8yrax9tdzw7jkecjxeg7st/programs/42f889e07ab07b4764f19207799046cb603b954659b601d1a1238aaeac111d5d.pl),consult/1)"
      """
