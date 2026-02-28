Feature: term_to_atom/2
  This feature is to test the term_to_atom/2 predicate.

  @great_for_documentation
  Scenario: Convert a ground term into a canonical atom
    This scenario demonstrates how `term_to_atom/2` turns a ground term into a canonical atom that can be reused later.

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      term_to_atom(greeting(hello, [world, 42]), Atom).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4172
      answer:
        has_more: false
        variables: ["Atom"]
        results:
        - substitutions:
          - variable: Atom
            expression: "'greeting(hello,[world,42])'"
      """

  @great_for_documentation
  Scenario: Parse a canonical atom back into a term
    This scenario demonstrates how `term_to_atom/2` reads an atom back into a Prolog term, including double-quoted strings.

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      term_to_atom(Term, 'payload(\"hi\", [foo, 42])').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5033
      answer:
        has_more: false
        variables: ["Term"]
        results:
        - substitutions:
          - variable: Term
            expression: payload([h,i],[foo,42])
      """

  Scenario: term_to_atom/2 requires one instantiated side

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      term_to_atom(Term, Atom).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3991
      answer:
        has_more: false
        variables: ["Term", "Atom"]
        results:
        - error: "error(instantiation_error,term_to_atom/2)"
      """

  Scenario: term_to_atom/2 rejects non-atom textual inputs

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      term_to_atom(Term, 42).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3984
      answer:
        has_more: false
        variables: ["Term"]
        results:
        - error: "error(type_error(atom,42),term_to_atom/2)"
      """

  Scenario: term_to_atom/2 raises a syntax_error on invalid canonical text

    Given the program:
      """ prolog
      """
    Given the query:
      """ prolog
      term_to_atom(Term, 'payload(').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4473
      answer:
        has_more: false
        variables: ["Term"]
        results:
        - error: "error(syntax_error(term),term_to_atom/2)"
      """
