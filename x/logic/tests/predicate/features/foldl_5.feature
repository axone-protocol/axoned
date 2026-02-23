Feature: foldl/5
  This feature is to test the foldl/5 predicate.

  @great_for_documentation
  Scenario: Fold two lists in lockstep to compute dot product
  This scenario demonstrates how to use foldl/5 to fold two lists simultaneously.

    Given the program:
      """ prolog
      add_product(X, Y, Acc0, Acc) :- Acc is Acc0 + (X * Y).
      """
    Given the query:
      """ prolog
      consult('/v1/lib/apply.pl'),
      foldl(add_product, [1,2,3], [4,5,6], 0, DotProduct).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3992
      answer:
        has_more: false
        variables: ["DotProduct"]
        results:
        - substitutions:
          - variable: DotProduct
            expression: 32
      """

  @great_for_documentation
  Scenario: Fold two empty lists returns the initial accumulator

    Given the program:
      """ prolog
      add_product(X, Y, Acc0, Acc) :- Acc is Acc0 + (X * Y).
      """
    Given the query:
      """ prolog
      consult('/v1/lib/apply.pl'),
      foldl(add_product, [], [], 99, Result).
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
            expression: 99
      """

  @great_for_documentation
  Scenario: Fold two lists to build a pair list

    Given the program:
      """ prolog
      make_pair(X, Y, Acc0, [[X,Y]|Acc0]).
      """
    Given the query:
      """ prolog
      consult('/v1/lib/apply.pl'),
      foldl(make_pair, [a,b], [1,2], [], Pairs).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3992
      answer:
        has_more: false
        variables: ["Pairs"]
        results:
        - substitutions:
          - variable: Pairs
            expression: "[[b,2],[a,1]]"
      """
