Feature: with_context/2
  This feature is to test the with_context/2 predicate.

  Scenario: with_context/2 succeeds when the wrapped goal succeeds
    This scenario demonstrates how with_context/2 transparently succeeds when the wrapped goal succeeds.

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      consult('/v1/lib/error.pl'),
      with_context(example/1, must_be(atom, hello)),
      Result = ok.
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4126
      answer:
        has_more: false
        variables: ["Result"]
        results:
        - substitutions:
          - variable: Result
            expression: ok
      """

  @great_for_documentation
  Scenario: with_context/2 rewrites the error context of the wrapped goal
    This scenario demonstrates how with_context/2 preserves the formal error while replacing its context.

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      consult('/v1/lib/error.pl'),
      with_context(example/1, must_be(atom, 42)).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4381
      answer:
        has_more: false
        results:
        - error: "error(type_error(atom,42),example/1)"
      """
