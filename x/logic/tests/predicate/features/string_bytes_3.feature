Feature: string_bytes/3
  This feature is to test the string_bytes/3 predicate.

  @great_for_documentation
  Scenario: Encode UTF-8 text into bytes
  This scenario demonstrates converting text into its UTF-8 byte sequence.

    Given the query:
      """ prolog
      string_bytes('aé', Bytes, utf8).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4392
      answer:
        has_more: false
        variables: ["Bytes"]
        results:
        - substitutions:
          - variable: Bytes
            expression: "[97,195,169]"
      """

  Scenario: Decode octet bytes into characters

    Given the query:
      """ prolog
      string_bytes(String, [249], octet).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4574
      answer:
        has_more: false
        variables: ["String"]
        results:
        - substitutions:
          - variable: String
            expression: "[ù]"
      """

  Scenario: Encode text with a specific charset

    Given the query:
      """ prolog
      string_bytes('a', Bytes, 'utf-16le').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5122
      answer:
        has_more: false
        variables: ["Bytes"]
        results:
        - substitutions:
          - variable: Bytes
            expression: "[97,0]"
      """

  Scenario: Reject an unknown charset

    Given the query:
      """ prolog
      string_bytes('a', Bytes, unknown).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5128
      answer:
        has_more: false
        variables: ["Bytes"]
        results:
        - error: "error(type_error(charset,unknown),string_bytes/3)"
      """

  Scenario: Require either text or bytes to be instantiated

    Given the query:
      """ prolog
      string_bytes(String, Bytes, utf8).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4827
      answer:
        has_more: false
        variables: ["String", "Bytes"]
        results:
        - error: "error(instantiation_error,string_bytes/3)"
      """
