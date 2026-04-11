Feature: base64url/2
  This feature is to test the base64url/2 predicate.

  @great_for_documentation
  Scenario: Encode and decode a string into a Base64 encoded atom in URL-Safe mode
  This scenario demonstrates how to encode an decode a plain string into a Base64-encoded atom using the `base64url/2`
  predicate.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64url('<<???>>', Encoded),
      base64url(Decoded, 'PDw_Pz8-Pg').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 15202
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

  Scenario: Error on incorrect URL-safe Base64 input
  This scenario demonstrates how `base64url/2` behaves when the encoded input is not valid URL-safe Base64 text.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64url(X, '!!!!').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 8387
      answer:
        has_more: false
        variables: ["X"]
        results:
        - error: "error(domain_error(encoding(base64),!!!!),base64_encoded/3)"
          substitutions:
      """

  Scenario: Error on non-canonical URL-safe Base64 tail
  This scenario demonstrates that `base64url/2` rejects malformed unpadded input when discarded tail bits are not zero.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64url(X, 'QR').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 7809
      answer:
        has_more: false
        variables: ["X"]
        results:
        - error: "error(domain_error(encoding(base64),QR),base64_encoded/3)"
          substitutions:
      """
