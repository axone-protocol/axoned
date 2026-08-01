Feature: hex_bytes/2
  This feature is to test the hex_bytes/2 predicate.

  @great_for_documentation
  Scenario: Decode a hexadecimal atom into bytes
    This scenario demonstrates how to decode a hexadecimal atom into a list of bytes.

    Given the query:
      """ prolog
      consult('/v1/lib/crypto.pl'),
      hex_bytes('501ACE', Bytes).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4929
      answer:
        has_more: false
        variables: ["Bytes"]
        results:
        - substitutions:
          - variable: Bytes
            expression: "[80,26,206]"
      """

  @great_for_documentation
  Scenario: Encode bytes into a hexadecimal atom
    This scenario demonstrates how to encode a list of bytes into a lowercase hexadecimal atom.

    Given the query:
      """ prolog
      consult('/v1/lib/crypto.pl'),
      hex_bytes(Hex, [80,26,206]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 8173
      answer:
        has_more: false
        variables: ["Hex"]
        results:
        - substitutions:
          - variable: Hex
            expression: "'501ace'"
      """

  Scenario: Reject hexadecimal character codes
    This scenario demonstrates that callers must explicitly convert character codes to an atom before calling hex_bytes/2.

    Given the query:
      """ prolog
      consult('/v1/lib/crypto.pl'),
      hex_bytes([53,48,49,65,67,69], Bytes).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4622
      answer:
        has_more: false
        variables: ["Bytes"]
        results:
        - error: "error(type_error(atom,[53,48,49,65,67,69]),hex_bytes/2)"
      """

  Scenario: Reject an invalid hexadecimal sequence
    This scenario demonstrates that malformed hexadecimal input raises a domain error.

    Given the query:
      """ prolog
      consult('/v1/lib/crypto.pl'),
      hex_bytes('501ACX', Bytes).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4845
      answer:
        has_more: false
        variables: ["Bytes"]
        results:
        - error: "error(domain_error(valid_encoding(hex),501ACX),hex_bytes/2)"
      """
