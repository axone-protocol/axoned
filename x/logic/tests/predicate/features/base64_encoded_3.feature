Feature: base64_encoded/3
  This feature is to test the base64_encoded/3 predicate.

  @great_for_documentation
  Scenario: Encode an atom into a Base64 encoded atom with default options
  This scenario demonstrates how to encode a plain atom into its Base64 representation using the `base64_encoded/3`
  predicate. The default options are used, meaning:
  - The output is returned as an atom.
  - Padding characters (`=`) are included (`padding(true)`).
  - The classic Base64 character set is used (`charset(classic)`), not the URL-safe variant.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded('Hello World', X, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 9406
      answer:
        has_more: false
        variables: ["X"]
        results:
        - substitutions:
          - variable: X
            expression: "'SGVsbG8gV29ybGQ='"
      """

  @great_for_documentation
  Scenario: Reject the removed output-shape option
    This scenario demonstrates that `base64_encoded/3` has a single textual output shape: atom.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded('Hello World', X, [as(atom)]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4407
      answer:
        has_more: false
        variables: ["X"]
        results:
        - error: "error(type_error(option,as(atom)),base64_encoded/3)"
          substitutions:
      """

  @great_for_documentation
  Scenario: Encode a string into a Base64 encoded atom without padding
  This scenario demonstrates how to encode a plain atom into a Base64-encoded atom using the `base64_encoded/3` predicate
  with custom options. The following options are used:
  - `padding(false)` – padding characters (`=`) are omitted.
  - The classic Base64 character set is used by default (`charset(classic)`).

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded('Hello World', X, [padding(false)]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 9982
      answer:
        has_more: false
        variables: ["X"]
        results:
        - substitutions:
          - variable: X
            expression: "'SGVsbG8gV29ybGQ'"
      """

  @great_for_documentation
  Scenario: Encode an atom into a Base64 encoded atom in URL-Safe mode
  This scenario demonstrates how to encode a plain atom into a Base64-encoded atom using the `base64_encoded/3` predicate
  with URL-safe encoding. The following options are used:
  - `charset(url)` – the URL-safe Base64 alphabet is used (e.g., `-` and `_` instead of `+` and `/`).
  - Padding characters are included by default (`padding(true)`).

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded('<<???>>', Classic, [charset(classic)]),
      base64_encoded('<<???>>', UrlSafe, [charset(url)]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 13251
      answer:
        has_more: false
        variables: ["Classic", "UrlSafe"]
        results:
        - substitutions:
          - variable: Classic
            expression: "'PDw/Pz8+Pg=='"
          - variable: UrlSafe
            expression: "'PDw_Pz8-Pg=='"
      """

  @great_for_documentation
  Scenario: Decode a Base64 encoded atom into a plain atom
  This scenario demonstrates how to decode a Base64-encoded atom back into plain text using the `base64_encoded/3` predicate.
  In this example, default options are used:
  •	The result is returned as an atom.
  •	Padding characters in the input are allowed (`padding(true)`).
  •	The classic Base64 character set is used (`charset(classic)`).

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded(X, 'SGVsbG8gV29ybGQ=', []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 11090
      answer:
        has_more: false
        variables: ["X"]
        results:
        - substitutions:
          - variable: X
            expression: "'Hello World'"
      """

  @great_for_documentation
  Scenario: Decode a Base64 encoded atom with explicit defaults
  This scenario demonstrates how to decode a Base64-encoded value back into plain text using the `base64_encoded/3` predicate.
  The following options are used:
  - `padding(true)` – padding characters in the input are allowed (default).
  - `charset(classic)` – the classic Base64 character set is used (default).

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded(X, 'SGVsbG8gV29ybGQ=', [padding(true), charset(classic)]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 12032
      answer:
        has_more: false
        variables: ["X"]
        results:
        - substitutions:
          - variable: X
            expression: "'Hello World'"
      """

  @great_for_documentation
  Scenario: Encode text using a specific character encoding
  This scenario demonstrates how the `encoding/1` option changes the bytes that are Base64-encoded before rendering the
  final Base64 text.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded('café', X, [encoding('iso-8859-1')]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 8279
      answer:
        has_more: false
        variables: ["X"]
        results:
        - substitutions:
          - variable: X
            expression: "'Y2Fm6Q=='"
      """

  Scenario: Reject a list of character codes
  This scenario demonstrates that callers must explicitly convert character codes to an atom before calling `base64_encoded/3`.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded([72,105], X, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5203
      answer:
        has_more: false
        variables: ["X"]
        results:
        - error: "error(type_error(atom,[72,105]),base64_encoded/3)"
          substitutions:
      """

  Scenario: Decode a Base64 encoded atom without padding
  This scenario demonstrates that decoding also supports `padding(false)` when the input omits trailing `=` characters.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded(X, 'SGVsbG8', [padding(false)]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 9115
      answer:
        has_more: false
        variables: ["X"]
        results:
        - substitutions:
          - variable: X
            expression: "'Hello'"
      """

  Scenario: Error when both arguments are variables
  This scenario demonstrates that `base64_encoded/3` requires at least one of `Plain` or `Encoded` to be instantiated.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded(X, Y, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4903
      answer:
        has_more: false
        variables: ["X", "Y"]
        results:
        - error: "error(instantiation_error,base64_encoded/3)"
          substitutions:
      """

  @great_for_documentation
  Scenario: Error on incorrect charset option
  This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid value is provided for the
  `charset` option.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded('Hello World', X, [charset(bad)]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4439
      answer:
        has_more: false
        variables: ["X"]
        results:
        - error: "error(domain_error(charset,bad),base64_encoded/3)"
          substitutions:
      """

  Scenario: Error on incorrect charset option (2)
  This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid type is provided for the
  `charset` option.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded('Hello World', X, [charset("bad")]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4647
      answer:
        has_more: false
        variables: ["X"]
        results:
        - error: "error(type_error(atom,[b,a,d]),base64_encoded/3)"
          substitutions:
      """

    @great_for_documentation
    Scenario: Error on incorrect padding option
    This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid value is provided for the
    `padding` option.

      Given the query:
        """ prolog
        consult('/v1/lib/base64.pl'),
        base64_encoded('Hello World', X, [padding(bad)]).
        """
      When the query is run
      Then the answer we get is:
        """ yaml
        height: 42
        gas_used: 4770
        answer:
          has_more: false
          variables: ["X"]
          results:
          - error: "error(domain_error(padding,bad),base64_encoded/3)"
            substitutions:
        """

  Scenario: Error on incorrect padding option (2)
  This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid type is provided for the
  `padding` option.

    Given the query:
        """ prolog
        consult('/v1/lib/base64.pl'),
        base64_encoded('Hello World', X, [padding(bad, 'very bad')]).
        """
    When the query is run
    Then the answer we get is:
        """ yaml
        height: 42
        gas_used: 4433
        answer:
          has_more: false
          variables: ["X"]
          results:
          - error: "error(type_error(option,padding(bad,very bad)),base64_encoded/3)"
            substitutions:
        """

  @great_for_documentation
  Scenario: Error on removed as option
  This scenario demonstrates how the `base64_encoded/3` predicate behaves when the removed `as` option is provided.

    Given the query:
        """ prolog
        consult('/v1/lib/base64.pl'),
        base64_encoded('Hello World', X, [as(bad)]).
        """
    When the query is run
    Then the answer we get is:
        """ yaml
        height: 42
        gas_used: 4406
        answer:
          has_more: false
          variables: ["X"]
          results:
          - error: "error(type_error(option,as(bad)),base64_encoded/3)"
            substitutions:
        """

  Scenario: Error on incorrect as option (2)
  This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid type is provided for the
  `as` option.

    Given the query:
        """ prolog
        consult('/v1/lib/base64.pl'),
        base64_encoded('Hello World', X, [as(bad, 'very bad')]).
        """
    When the query is run
    Then the answer we get is:
        """ yaml
        height: 42
        gas_used: 4428
        answer:
          has_more: false
          variables: ["X"]
          results:
          - error: "error(type_error(option,as(bad,very bad)),base64_encoded/3)"
            substitutions:
        """

  Scenario: Error on incorrect plain type input
  This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid plain type input is provided.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded(wrong(input), X, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5195
      answer:
        has_more: false
        variables: ["X"]
        results:
        - error: "error(type_error(atom,wrong(input)),base64_encoded/3)"
          substitutions:
      """

  Scenario: Error on incorrect Base64 encoded input
  This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid encoded input is provided.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded(X, '!!!!', []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5805
      answer:
        has_more: false
        variables: ["X"]
        results:
        - error: "error(domain_error(encoding(base64),!!!!),base64_encoded/3)"
          substitutions:
      """

  Scenario: Error on non-canonical 2-character tail without padding
  This scenario demonstrates that `base64_encoded/3` rejects a malformed final quantum when discarded tail bits are not zero.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded(X, 'QR', [padding(false)]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 6329
      answer:
        has_more: false
        variables: ["X"]
        results:
        - error: "error(domain_error(encoding(base64),QR),base64_encoded/3)"
          substitutions:
      """

  Scenario: Error on non-canonical 3-character tail without padding
  This scenario demonstrates that `base64_encoded/3` rejects a malformed 3-character tail when discarded tail bits are not zero.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded(X, 'SGl', [padding(false)]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 6890
      answer:
        has_more: false
        variables: ["X"]
        results:
        - error: "error(domain_error(encoding(base64),SGl),base64_encoded/3)"
          substitutions:
      """

  Scenario: Error on extra data after padded Base64 input
  This scenario demonstrates that `base64_encoded/3` rejects additional data after a padded final quantum.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded(X, 'QQ==QQ==', []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 7028
      answer:
        has_more: false
        variables: ["X"]
        results:
        - error: "error(domain_error(encoding(base64),QQ==QQ==),base64_encoded/3)"
          substitutions:
      """

  Scenario: Error on incorrect encoded type input
  This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid encoded type input is provided.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded(X, wrong(input), []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5331
      answer:
        has_more: false
        variables: ["X"]
        results:
        - error: "error(type_error(atom,wrong(input)),base64_encoded/3)"
          substitutions:
      """

  @great_for_documentation
  Scenario: Error on incorrect encoding option
  This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid value is provided for the
  `encoding` option.

    Given the query:
        """ prolog
        consult('/v1/lib/base64.pl'),
        base64_encoded(X, 'SGVsbG8gV29ybGQ=', [encoding(unknown)]).
        """
    When the query is run
    Then the answer we get is:
        """ yaml
        height: 42
        gas_used: 11778
        answer:
          has_more: false
          variables: ["X"]
          results:
          - error: "error(type_error(charset,unknown),base64_encoded/3)"
            substitutions:
        """

  @great_for_documentation
  Scenario: Error on incorrect encoding option (2)
  This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid type is provided for the
  `encoding` option.

    Given the query:
        """ prolog
        consult('/v1/lib/base64.pl'),
        base64_encoded(X, 'SGVsbG8gV29ybGQ=', [encoding(bad, 'very bad')]).
        """
    When the query is run
    Then the answer we get is:
        """ yaml
        height: 42
        gas_used: 4439
        answer:
          has_more: false
          variables: ["X"]
          results:
          - error: "error(type_error(option,encoding(bad,very bad)),base64_encoded/3)"
            substitutions:
        """

  @great_for_documentation
  Scenario: Error on unknown option name
  This scenario demonstrates how the `base64_encoded/3` predicate behaves when an unknown option name is provided.
  This helps catch typos in option names (e.g., `chatset` instead of `charset`).

    Given the query:
        """ prolog
        consult('/v1/lib/base64.pl'),
        base64_encoded('Hello World', X, [chatset(classic)]).
        """
    When the query is run
    Then the answer we get is:
        """ yaml
        height: 42
        gas_used: 4415
        answer:
          has_more: false
          variables: ["X"]
          results:
          - error: "error(type_error(option,chatset(classic)),base64_encoded/3)"
            substitutions:
        """
