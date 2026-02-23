Feature: foldl/7
  This feature is to test the foldl/7 predicate.

  @great_for_documentation
  Scenario: Fold four lists in lockstep to compute a combined result
  This scenario demonstrates how to use foldl/7 to fold four lists simultaneously.

    Given the program:
      """ prolog
      quad_sum(W, X, Y, Z, Acc0, Acc) :- Acc is Acc0 + W + X + Y + Z.
      """
    Given the query:
      """ prolog
      consult('/v1/lib/apply.pl'),
      foldl(quad_sum, [1,2], [3,4], [5,6], [7,8], 0, Result).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3992
      answer:
        has_more: false
        variables: ["Result"]
        results:
        - substitutions:
          - variable: Result
            expression: 36
      """

  @great_for_documentation
  Scenario: Fold four empty lists returns the initial accumulator

    Given the program:
      """ prolog
      quad_sum(W, X, Y, Z, Acc0, Acc) :- Acc is Acc0 + W + X + Y + Z.
      """
    Given the query:
      """ prolog
      consult('/v1/lib/apply.pl'),
      foldl(quad_sum, [], [], [], [], 100, Result).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3992
      answer:
        has_more: false
        variables: ["Result"]
        results:
        - substitutions:
          - variable: Result
            expression: 100
      """

  @great_for_documentation
  Scenario: Fold four lists to build a structured result

    Given the program:
      """ prolog
      make_quad(W, X, Y, Z, Acc0, [[W,X,Y,Z]|Acc0]).
      """
    Given the query:
      """ prolog
      consult('/v1/lib/apply.pl'),
      foldl(make_quad, [a], [1], [x], [true], [], Quads).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3992
      answer:
        has_more: false
        variables: ["Quads"]
        results:
        - substitutions:
          - variable: Quads
            expression: "[[a,1,x,true]]"
      """
