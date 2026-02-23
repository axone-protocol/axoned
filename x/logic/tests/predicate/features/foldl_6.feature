Feature: foldl/6
  This feature is to test the foldl/6 predicate.

  @great_for_documentation
  Scenario: Fold three lists in lockstep to compute weighted sum
  This scenario demonstrates how to use foldl/6 to fold three lists simultaneously.

    Given the program:
      """ prolog
      weighted_sum(X, Y, Z, Acc0, Acc) :- Acc is Acc0 + (X * Y * Z).
      """
    Given the query:
      """ prolog
      consult('/v1/lib/apply.pl'),
      foldl(weighted_sum, [1,2], [3,4], [5,6], 0, Result).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3984
      answer:
        has_more: false
        variables: ["Result"]
        results:
        - substitutions:
          - variable: Result
            expression: 63
      """

  @great_for_documentation
  Scenario: Fold three empty lists returns the initial accumulator

    Given the program:
      """ prolog
      weighted_sum(X, Y, Z, Acc0, Acc) :- Acc is Acc0 + (X * Y * Z).
      """
    Given the query:
      """ prolog
      consult('/v1/lib/apply.pl'),
      foldl(weighted_sum, [], [], [], 42, Result).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3976
      answer:
        has_more: false
        variables: ["Result"]
        results:
        - substitutions:
          - variable: Result
            expression: 42
      """

  @great_for_documentation
  Scenario: Fold three lists to build a structured result

    Given the program:
      """ prolog
      make_triple(X, Y, Z, Acc0, [[X,Y,Z]|Acc0]).
      """
    Given the query:
      """ prolog
      consult('/v1/lib/apply.pl'),
      foldl(make_triple, [a,b], [1,2], [x,y], [], Triples).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3982
      answer:
        has_more: false
        variables: ["Triples"]
        results:
        - substitutions:
          - variable: Triples
            expression: "[[b,2,y],[a,1,x]]"
      """
