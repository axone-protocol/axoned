Feature: foldl/4
  This feature is to test the foldl/4 predicate.

  @great_for_documentation
  Scenario: Fold a list of integers into a sum
  This scenario demonstrates how to load apply.pl and use foldl/4 to aggregate a list with an accumulator.

    Given the program:
      """ prolog
      sum(Elem, Acc0, Acc) :- Acc is Acc0 + Elem.
      """
    Given the query:
      """ prolog
      consult('/v1/lib/apply.pl'),
      foldl(sum, [1,2,3,4], 0, Total).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3992
      answer:
        has_more: false
        variables: ["Total"]
        results:
        - substitutions:
          - variable: Total
            expression: 10
      """

  Scenario: foldl/4 is unavailable until apply.pl is loaded

    Given the program:
      """ prolog
      sum(Elem, Acc0, Acc) :- Acc is Acc0 + Elem.
      """
    Given the query:
      """ prolog
      foldl(sum, [1,2], 0, Total).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3975
      answer:
        has_more: false
        variables: ["Total"]
        results:
        - error: "error(existence_error(procedure,foldl/4),root)"
      """
