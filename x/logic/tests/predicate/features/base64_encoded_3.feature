Feature: base64_encoded/3
  This feature is to test the base64_encoded/3 predicate.

  @great_for_documentation
  Scenario: Encode a string into a Base64 encoded string (with default options)
  This scenario demonstrates how to encode a plain string into its Base64 representation using the `base64_encoded/3`
  predicate. The default options are used, meaning:
  - The output is returned as a list of characters (`as(string)`).
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
      gas_used: 9043
      answer:
        has_more: false
        variables: ["X"]
        results:
        - substitutions:
          - variable: X
            expression: "['S','G','V',s,b,'G','8',g,'V','2','9',y,b,'G','Q',=]"
      """

  @great_for_documentation
  Scenario: Encode a string into a Base64 encoded atom
    This scenario demonstrates how to encode a plain string into a Base64-encoded atom using the `base64_encoded/3`
    predicate. The `as(atom)` option is specified, so the result is returned as a Prolog atom instead of a character
    list. All other options use their default values.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded('Hello World', X, [as(atom)]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 9408
      answer:
        has_more: false
        variables: ["X"]
        results:
        - substitutions:
          - variable: X
            expression: "'SGVsbG8gV29ybGQ='"
      """

  @great_for_documentation
  Scenario: Encode a string into a Base64 encoded atom without padding
  This scenario demonstrates how to encode a plain string into a Base64-encoded atom using the `base64_encoded/3` predicate
  with custom options. The following options are used:
  - `as(atom)` – the result is returned as a Prolog atom.
  - `padding(false)` – padding characters (`=`) are omitted.
  - The classic Base64 character set is used by default (`charset(classic)`).

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded('Hello World', X, [as(atom), padding(false)]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 10063
      answer:
        has_more: false
        variables: ["X"]
        results:
        - substitutions:
          - variable: X
            expression: "'SGVsbG8gV29ybGQ'"
      """

  @great_for_documentation
  Scenario: Encode a String into a Base64 encoded atom in URL-Safe mode
  This scenario demonstrates how to encode a plain string into a Base64-encoded atom using the `base64_encoded/3` predicate
  with URL-safe encoding. The following options are used:
  - `as(atom)` – the result is returned as a Prolog atom.
  - `charset(url)` – the URL-safe Base64 alphabet is used (e.g., `-` and `_` instead of `+` and `/`).
  - Padding characters are included by default (`padding(true)`).

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded('<<???>>', Classic, [as(atom), charset(classic)]),
      base64_encoded('<<???>>', UrlSafe, [as(atom), charset(url)]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 13901
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
  Scenario: Decode a Base64 encoded String into plain text
  This scenario demonstrates how to decode a Base64-encoded value back into plain text using the `base64_encoded/3` predicate.
  The encoded input can be provided as a character list or an atom. In this example, default options are used:
  •	The result (plain text) is returned as a character list (`as(string)`).
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
      gas_used: 10809
      answer:
        has_more: false
        variables: ["X"]
        results:
        - substitutions:
          - variable: X
            expression: "['H',e,l,l,o,' ','W',o,r,l,d]"
      """

  @great_for_documentation
  Scenario: Decode a Base64 Encoded string into a plain atom
  This scenario demonstrates how to decode a Base64-encoded value back into plain text using the `base64_encoded/3` predicate,
  with the result returned as a Prolog atom. The following options are used:
  - `as(atom)` – the decoded plain text is returned as an atom.
  - `padding(true)` – padding characters in the input are allowed (default).
  - `charset(classic)` – the classic Base64 character set is used (default).

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded(X, 'SGVsbG8gV29ybGQ=', [as(atom)]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 11678
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
      base64_encoded('café', X, [as(atom), encoding('iso-8859-1')]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 8016
      answer:
        has_more: false
        variables: ["X"]
        results:
        - substitutions:
          - variable: X
            expression: "'Y2Fm6Q=='"
      """

  Scenario: Encode a list of character codes
  This scenario demonstrates that `base64_encoded/3` accepts a list of character codes as plain text input.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded([72,105], X, [as(atom)]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 6595
      answer:
        has_more: false
        variables: ["X"]
        results:
        - substitutions:
          - variable: X
            expression: "'SGk='"
      """

  Scenario: Decode a Base64 encoded atom without padding
  This scenario demonstrates that decoding also supports `padding(false)` when the input omits trailing `=` characters.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded(X, 'SGVsbG8', [as(atom), padding(false)]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 8543
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
      gas_used: 5099
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
      gas_used: 4445
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
      gas_used: 4653
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
        gas_used: 4801
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
        gas_used: 4469
        answer:
          has_more: false
          variables: ["X"]
          results:
          - error: "error(type_error(option,padding(bad,very bad)),base64_encoded/3)"
            substitutions:
        """

  @great_for_documentation
  Scenario: Error on incorrect as option
  This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid value is provided for the
  `as` option.

    Given the query:
        """ prolog
        consult('/v1/lib/base64.pl'),
        base64_encoded('Hello World', X, [as(bad)]).
        """
    When the query is run
    Then the answer we get is:
        """ yaml
        height: 42
        gas_used: 5253
        answer:
          has_more: false
          variables: ["X"]
          results:
          - error: "error(domain_error(as,bad),base64_encoded/3)"
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
        gas_used: 4464
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
      gas_used: 5070
      answer:
        has_more: false
        variables: ["X"]
        results:
        - error: "error(type_error(text,wrong(input)),base64_encoded/3)"
          substitutions:
      """

  Scenario: Error on incorrect Base64 encoded input
  This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid encoded input is provided.

    Given the query:
      """ prolog
      consult('/v1/lib/base64.pl'),
      base64_encoded(X, '!!!!', [as(atom)]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 6537
      answer:
        has_more: false
        variables: ["X"]
        results:
        - error: "error(domain_error(encoding(base64),!!!!),base64_encoded/3)"
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
      gas_used: 5200
      answer:
        has_more: false
        variables: ["X"]
        results:
        - error: "error(type_error(text,wrong(input)),base64_encoded/3)"
          substitutions:
      """

  @great_for_documentation
  Scenario: Error on incorrect encoding option
  This scenario demonstrates how the `base64_encoded/3` predicate behaves when an invalid value is provided for the
  `encoding` option.

    Given the query:
        """ prolog
        consult('/v1/lib/base64.pl'),
        base64_encoded(X, 'SGVsbG8gV29ybGQ=', [as(atom), encoding(unknown)]).
        """
    When the query is run
    Then the answer we get is:
        """ yaml
        height: 42
        gas_used: 12449
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
        gas_used: 4475
        answer:
          has_more: false
          variables: ["X"]
          results:
          - error: "error(type_error(option,encoding(bad,very bad)),base64_encoded/3)"
            substitutions:
        """
