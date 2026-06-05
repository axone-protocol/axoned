Feature: source_file/1
  This feature is to test the source_file/1 predicate.

  @great_for_documentation
  Scenario: Match a loaded source file
  This scenario demonstrates checking whether a source file has been loaded.

    Given the query:
      """ prolog
      consult('/v1/lib/lists.pl'),
      source_file('/v1/lib/lists.pl').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4301
      answer:
        has_more: false
        variables:
        results:
        - substitutions:
      """

  Scenario: Enumerate loaded source files

    Given the query:
      """ prolog
      consult('/v1/lib/lists.pl'),
      source_file(File).
      """
    When the query is run (limited to 1 solutions)
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4215
      answer:
        has_more: false
        variables: ["File"]
        results:
        - substitutions:
          - variable: File
            expression: "'/v1/lib/lists.pl'"
      """

  Scenario: Fail when the source file is unknown

    Given the query:
      """ prolog
      source_file('/v1/lib/missing.pl').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4227
      answer:
        has_more: false
        variables:
        results:
      """

  Scenario: Reject non-atom source file input

    Given the query:
      """ prolog
      source_file(42).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4086
      answer:
        has_more: false
        results:
        - error: "error(type_error(atom,42),source_file/1)"
      """
