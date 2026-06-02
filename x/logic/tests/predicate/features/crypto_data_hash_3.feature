Feature: crypto_data_hash/3
  This feature is to test the crypto_data_hash/3 predicate.

  @great_for_documentation
  Scenario: Compute a SHA-256 hash with default options
    This scenario demonstrates how to compute a SHA-256 digest from text using the default options.

    Given the query:
      """ prolog
      consult('/v1/lib/crypto.pl'),
      crypto_data_hash('hello world', Hash, []),
      hex_bytes(Hex, Hash).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 64284
      answer:
        has_more: false
        variables: ["Hash", "Hex"]
        results:
        - substitutions:
          - variable: Hash
            expression: "[185,77,39,185,147,77,62,8,165,46,82,215,218,125,171,250,196,132,239,227,122,83,128,238,144,136,247,172,226,239,205,233]"
          - variable: Hex
            expression: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
      """

  Scenario: Compute hashes with explicit algorithms
    This scenario demonstrates how to select SHA-512 and MD5 explicitly.

    Given the query:
      """ prolog
      consult('/v1/lib/crypto.pl'),
      crypto_data_hash('hello world', Sha512Bytes, [algorithm(sha512)]),
      crypto_data_hash('hello world', Md5Bytes, [algorithm(md5)]),
      hex_bytes(Sha512, Sha512Bytes),
      hex_bytes(Md5, Md5Bytes).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 144730
      answer:
        has_more: false
        variables: ["Sha512Bytes", "Md5Bytes", "Sha512", "Md5"]
        results:
        - substitutions:
          - variable: Sha512Bytes
            expression: "[48,158,204,72,156,18,214,235,76,196,15,80,201,2,242,180,208,237,119,238,81,26,124,122,155,205,60,168,109,76,216,111,152,157,211,91,197,255,73,150,112,218,52,37,91,69,176,207,216,48,232,31,96,93,207,125,197,84,46,147,174,156,215,111]"
          - variable: Md5Bytes
            expression: "[94,182,59,187,224,30,238,208,147,203,34,187,143,90,205,195]"
          - variable: Sha512
            expression: "'309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f'"
          - variable: Md5
            expression: "'5eb63bbbe01eeed093cb22bb8f5acdc3'"
      """

  Scenario: Compute hashes from hexadecimal and octet data
    This scenario demonstrates how to hash data already represented as bytes.

    Given the query:
      """ prolog
      consult('/v1/lib/crypto.pl'),
      crypto_data_hash('68656c6c6f20776f726c64', HashFromHex, [encoding(hex)]),
      crypto_data_hash([104,101,108,108,111,32,119,111,114,108,100], HashFromOctet, [encoding(octet)]),
      hex_bytes(HexFromHex, HashFromHex),
      hex_bytes(HexFromOctet, HashFromOctet).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 129902
      answer:
        has_more: false
        variables: ["HashFromHex", "HashFromOctet", "HexFromHex", "HexFromOctet"]
        results:
        - substitutions:
          - variable: HashFromHex
            expression: "[185,77,39,185,147,77,62,8,165,46,82,215,218,125,171,250,196,132,239,227,122,83,128,238,144,136,247,172,226,239,205,233]"
          - variable: HashFromOctet
            expression: "[185,77,39,185,147,77,62,8,165,46,82,215,218,125,171,250,196,132,239,227,122,83,128,238,144,136,247,172,226,239,205,233]"
          - variable: HexFromHex
            expression: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
          - variable: HexFromOctet
            expression: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
      """

  Scenario: Match a computed hash against an expected value
    This scenario demonstrates that crypto_data_hash/3 can verify a provided digest.

    Given the query:
      """ prolog
      consult('/v1/lib/crypto.pl'),
      hex_bytes('b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9', Expected),
      crypto_data_hash('hello world', Expected, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 16965
      answer:
        has_more: false
        variables: ["Expected"]
        results:
        - substitutions:
          - variable: Expected
            expression: "[185,77,39,185,147,77,62,8,165,46,82,215,218,125,171,250,196,132,239,227,122,83,128,238,144,136,247,172,226,239,205,233]"
      """

  Scenario: Reject an unknown hash algorithm
    This scenario demonstrates that unsupported algorithms raise a type error.

    Given the query:
      """ prolog
      consult('/v1/lib/crypto.pl'),
      crypto_data_hash('hello world', Hash, [algorithm(cheh)]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4573
      answer:
        has_more: false
        variables: ["Hash"]
        results:
        - error: "error(type_error(hash_algorithm,cheh),crypto_data_hash/3)"
      """

  Scenario: Reject an unbound option term
    This scenario demonstrates that an unbound option raises an instantiation error.

    Given the query:
      """ prolog
      consult('/v1/lib/crypto.pl'),
      crypto_data_hash('hello world', Hash, [Opt]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4303
      answer:
        has_more: false
        variables: ["Hash", "Opt"]
        results:
        - error: "error(instantiation_error,crypto_data_hash/3)"
      """
