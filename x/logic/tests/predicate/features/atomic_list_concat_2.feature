Feature: atomic_list_concat/2
  This feature is to test the atomic_list_concat/2 predicate.

  @great_for_documentation
  Scenario: Concatenate atomic values into a single atom
    This scenario demonstrates how `atomic_list_concat/2` concatenates the textual representation of several atomic values.

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      atomic_list_concat([hello, '-', 42, '-', world], Atom).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4295
      answer:
        has_more: false
        variables: ["Atom"]
        results:
        - substitutions:
          - variable: Atom
            expression: "'hello-42-world'"
      """

  Scenario: atomic_list_concat/2 requires an instantiated list

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      atomic_list_concat(List, Atom).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3981
      answer:
        has_more: false
        variables: ["List", "Atom"]
        results:
        - error: "error(instantiation_error,atomic_list_concat/2)"
      """

  Scenario: atomic_list_concat/2 rejects partial lists

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      List = [hello|Tail],
      atomic_list_concat(List, Atom).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3987
      answer:
        has_more: false
        variables: ["List", "Tail", "Atom"]
        results:
        - error: "error(instantiation_error,atomic_list_concat/2)"
      """

  Scenario: atomic_list_concat/2 rejects non-list values

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      atomic_list_concat(foo, Atom).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3989
      answer:
        has_more: false
        variables: ["Atom"]
        results:
        - error: "error(type_error(list,foo),atomic_list_concat/2)"
      """
