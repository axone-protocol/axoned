Feature: eddsa_verify/4
  This feature is to test the eddsa_verify/4 predicate.

  @great_for_documentation
  Scenario: Verify an Ed25519 signature with default options
    This scenario demonstrates how to verify an Ed25519 signature over hexadecimal data.

    Given the program:
      """ prolog
      valid_ed25519(Verified) :-
        hex_bytes('53167ac3fc4b720daa45b04fc73fe752578fa23a10048422d6904b7f4f7bba5a', PubKey),
        hex_bytes('889bcfd331e8e43b5ebf430301dffb6ac9e2fce69f6227b43552fe3dc8cc1ee00c1cc53452a8712e9d5f80086dff8cf4999c1b93ed6c6e403c09334cb61ddd0b', Signature),
        eddsa_verify(PubKey, '9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', Signature, []),
        Verified = true.
      """
    Given the query:
      """ prolog
      consult('/v1/lib/crypto.pl'),
      valid_ed25519(Verified).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 235323
      answer:
        has_more: false
        variables: ["Verified"]
        results:
        - substitutions:
          - variable: Verified
            expression: true
      """

  Scenario: Verify an Ed25519 signature over octet data
    This scenario demonstrates how to verify an Ed25519 signature when the data is already a list of bytes.

    Given the program:
      """ prolog
      valid_ed25519_octet(Verified) :-
        hex_bytes('53167ac3fc4b720daa45b04fc73fe752578fa23a10048422d6904b7f4f7bba5a', PubKey),
        hex_bytes('9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', Data),
        hex_bytes('889bcfd331e8e43b5ebf430301dffb6ac9e2fce69f6227b43552fe3dc8cc1ee00c1cc53452a8712e9d5f80086dff8cf4999c1b93ed6c6e403c09334cb61ddd0b', Signature),
        eddsa_verify(PubKey, Data, Signature, [encoding(octet), type(ed25519)]),
        Verified = true.
      """
    Given the query:
      """ prolog
      consult('/v1/lib/crypto.pl'),
      valid_ed25519_octet(Verified).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 241727
      answer:
        has_more: false
        variables: ["Verified"]
        results:
        - substitutions:
          - variable: Verified
            expression: true
      """

  @great_for_documentation
  Scenario: Reject an invalid Ed25519 signature
    This scenario demonstrates that eddsa_verify/4 fails when the signature does not match the data.

    Given the program:
      """ prolog
      invalid_ed25519 :-
        hex_bytes('53167ac3fc4b720daa45b04fc73fe752578fa23a10048422d6904b7f4f7bba5a', PubKey),
        hex_bytes('889bcfd331e8e43b5ebf430301dffb6ac9e2fce69f6227b43552fe3dc8cc1ee00c1cc53452a8712e9d5f80086dff8cf4999c1b93ed6c6e403c09334cb61ddd0b', Signature),
        eddsa_verify(PubKey, '9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9e', Signature, []).
      """
    Given the query:
      """ prolog
      consult('/v1/lib/crypto.pl'),
      invalid_ed25519.
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 235291
      answer:
        has_more: false
        variables:
        results:
      """

  @great_for_documentation
  Scenario: Reject an unsupported EdDSA algorithm
    This scenario demonstrates that eddsa_verify/4 rejects algorithms outside the EdDSA family.

    Given the program:
      """ prolog
      unsupported_eddsa :-
        eddsa_verify([], '', [], [type(secp256k1)]).
      """
    Given the query:
      """ prolog
      consult('/v1/lib/crypto.pl'),
      unsupported_eddsa.
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4809
      answer:
        has_more: false
        variables:
        results:
        - error: "error(type_error(cryptographic_algorithm,secp256k1),eddsa_verify/4)"
      """

  Scenario: Reject malformed EdDSA option terms
    This scenario demonstrates that eddsa_verify/4 rejects malformed option terms instead of falling back to defaults.

    Given the program:
      """ prolog
      malformed_eddsa_option :-
        eddsa_verify([], '', [], [encoding(hex, utf8)]).
      """
    Given the query:
      """ prolog
      consult('/v1/lib/crypto.pl'),
      malformed_eddsa_option.
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4567
      answer:
        has_more: false
        variables:
        results:
        - error: "error(type_error(option,encoding(hex,utf8)),eddsa_verify/4)"
      """
