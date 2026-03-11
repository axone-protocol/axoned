Feature: open/3
  This feature is to test the open/3 predicate.

  @great_for_documentation
  Scenario: Open a snapshot resource for reading
  This scenario demonstrates how to open a read-only snapshot resource exposed by the VFS.

    Given the query:
      """ prolog
      open('/v1/run/header/height', read, _).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3953
      answer:
        has_more: false
        variables:
        results:
        - substitutions:
      """
