Feature: base64/2
  This feature is to test the base64/2 predicate.

  @great_for_documentation
  Scenario: Encode and decode a string into a Base64 encoded atom
  This scenario demonstrates how to encode an decode a plain string into a Base64-encoded atom using the `base64/2`
  predicate.

    Given the query:
      """ prolog
      base64('Hello world', Encoded),
      base64(Decoded, 'SGVsbG8gd29ybGQ=').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3976
      answer:
        has_more: false
        variables: ["Encoded", "Decoded"]
        results:
        - substitutions:
          - variable: Encoded
            expression: "'SGVsbG8gd29ybGQ='"
          - variable: Decoded
            expression: "'Hello world'"
      """
