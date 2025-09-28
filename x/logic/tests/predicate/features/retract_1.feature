Feature: retract/1
  This feature is to test the retract/1 predicate.

  @great_for_documentation
  Scenario: Retract a fact from the database.
    This scenario demonstrates the process of retracting a fact from a Prolog database. In Prolog, retracting a fact means
    removing a piece of information or *knowledge* from the database, making it unavailable for subsequent queries.
    This is particularly useful when you want to dynamically remove facts or rules based on conditions or interactions
    during runtime.

    Given the query:
      """ prolog
      assertz(parent(john, alice)), retract(parent(john, alice)), parent(X, Y).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3977
      answer:
        has_more: false
        variables: ["X","Y"]
        results:
      """

  @great_for_documentation
  Scenario: Only dynamic predicates can be retracted.
    This scenario demonstrates that only dynamic predicates can be retracted. In Prolog, dynamic predicates are those that can be
    modified during runtime. This is in contrast to static predicates, which are fixed and cannot be modified.

    Given the program:
      """ prolog
      parent(jane, alice).
      """
    Given the query:
      """ prolog
      retract(parent(jane, alice)).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3975
      answer:
        has_more: false
        results:
        - error: "error(permission_error(modify,static_procedure,parent/2),retract/1)"
      """
