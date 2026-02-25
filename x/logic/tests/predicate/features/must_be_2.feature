Feature: must_be/2
  This feature is to test the must_be/2 predicate.

  Scenario: must_be/2 is unavailable until error.pl is loaded

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      must_be(atom, hello).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3975
      answer:
        has_more: false
        results:
        - error: "error(existence_error(procedure,must_be/2),root)"
      """

  @great_for_documentation
  Scenario: Validate an atom with must_be/2
    This scenario demonstrates how to load error.pl and validate a value type with must_be/2.

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      consult('/v1/lib/error.pl'),
      must_be(atom, hello),
      Result = ok.
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3986
      answer:
        has_more: false
        variables: ["Result"]
        results:
        - substitutions:
          - variable: Result
            expression: ok
      """

  @great_for_documentation
  Scenario: must_be/2 throws instantiation_error for unbound values
    This scenario demonstrates that must_be/2 raises an instantiation error when the checked value is a variable.

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      consult('/v1/lib/error.pl'),
      must_be(atom, X).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3997
      answer:
        has_more: false
        variables: ["X"]
        results:
        - error: "error(instantiation_error,must_be/2)"
      """

  Scenario: must_be/2 throws type_error when type does not match

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      consult('/v1/lib/error.pl'),
      must_be(integer, hello).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3997
      answer:
        has_more: false
        results:
        - error: "error(type_error(integer,hello),must_be/2)"
      """

  Scenario: must_be/2 throws instantiation_error on partial list

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      consult('/v1/lib/error.pl'),
      Partial = [a|Tail],
      must_be(list, Partial).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4005
      answer:
        has_more: false
        variables: ["Partial", "Tail"]
        results:
        - error: "error(instantiation_error,must_be/2)"
      """
