Feature: atomic_list_concat/3
  This feature is to test the atomic_list_concat/3 predicate.

  @great_for_documentation
  Scenario: Concatenate values with a separator
    This scenario demonstrates how `atomic_list_concat/3` inserts a separator between the textual representation of each list element.

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      atomic_list_concat([cosmos, hub, 4], '-', Atom).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4118
      answer:
        has_more: false
        variables: ["Atom"]
        results:
        - substitutions:
          - variable: Atom
            expression: "'cosmos-hub-4'"
      """

  @great_for_documentation
  Scenario: Build a URI-like atom from separate parts
    This scenario demonstrates how `atomic_list_concat/3` can be used to assemble a structured atom from reusable parts.

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      atomic_list_concat([scheme, host, path], '://', URI).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4150
      answer:
        has_more: false
        variables: ["URI"]
        results:
        - substitutions:
          - variable: URI
            expression: "'scheme://host://path'"
      """

  Scenario: atomic_list_concat/3 requires an instantiated separator

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      atomic_list_concat([a, b], Separator, Atom).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4012
      answer:
        has_more: false
        variables: ["Separator", "Atom"]
        results:
        - error: "error(instantiation_error,atomic_list_concat/3)"
      """

  Scenario: atomic_list_concat/3 rejects non-atom separators

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      atomic_list_concat([a, b], 42, Atom).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4016
      answer:
        has_more: false
        variables: ["Atom"]
        results:
        - error: "error(type_error(atom,42),atomic_list_concat/3)"
      """
