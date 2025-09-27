Feature: base64url/2
  This feature is to test the base64url/2 predicate.

  @great_for_documentation
  Scenario: Encode and decode a string into a Base64 encoded atom in URL-Safe mode
  This scenario demonstrates how to encode an decode a plain string into a Base64-encoded atom using the `base64url/2`
  predicate.

    Given the query:
      """ prolog
      base64url('<<???>>', Encoded),
      base64url(Decoded, 'PDw_Pz8-Pg').
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
            expression: "'PDw_Pz8-Pg'"
          - variable: Decoded
            expression: "<<???>>"
      """
