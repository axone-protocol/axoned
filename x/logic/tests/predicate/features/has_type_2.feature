Feature: has_type/2
  This feature is to test the has_type/2 predicate.

  Scenario: has_type/2 is unavailable until type.pl is loaded

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      has_type(atom, hello).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3917
      answer:
        has_more: false
        results:
        - error: "error(existence_error(procedure,has_type/2),root)"
      """

  @great_for_documentation
  Scenario: Validate a byte list with has_type/2
    This scenario demonstrates how to load type.pl and check a structured type using has_type/2.

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      consult('/v1/lib/type.pl'),
      has_type(list(byte), [0,1,255]),
      Result = ok.
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4488
      answer:
        has_more: false
        variables: ["Result"]
        results:
        - substitutions:
          - variable: Result
            expression: ok
      """

  @great_for_documentation
  Scenario: has_type/2 fails when the type does not match
    This scenario demonstrates that has_type/2 fails quietly when the value does not satisfy the requested type.

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      consult('/v1/lib/type.pl'),
      has_type(integer, hello).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4019
      answer:
        has_more: false
        results:
      """
