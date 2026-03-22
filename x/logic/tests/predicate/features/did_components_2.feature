Feature: did_components/2
  This feature is to test the did_components/2 predicate.

  @great_for_documentation
  Scenario: Parse a DID URL into raw DID components
    This scenario demonstrates how to decompose a DID URL into a `did/5` structured term.
    Path is preserved with its leading `/`, while query and fragment are preserved raw without their leading separators.

    Given the query:
      """ prolog
      consult('/v1/lib/did.pl'),
      did_components(
        'did:example:123456/path?versionId=1#auth-key',
        did(Method, MethodSpecificId, Path, Query, Fragment)
      ).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 11918
      answer:
        has_more: false
        variables: ["Method", "MethodSpecificId", "Path", "Query", "Fragment"]
        results:
        - substitutions:
          - variable: Method
            expression: "example"
          - variable: MethodSpecificId
            expression: "'123456'"
          - variable: Path
            expression: "'/path'"
          - variable: Query
            expression: "'versionId=1'"
          - variable: Fragment
            expression: "'auth-key'"
      """

  @great_for_documentation
  Scenario: Reconstruct a DID URL from raw DID components
    This scenario demonstrates the reverse mode of did_components/2, reconstructing a DID URL from a `did/5` term.

    Given the query:
      """ prolog
      consult('/v1/lib/did.pl'),
      did_components(DID, did(example, '123456', '/foo/bar', 'versionId=1', test)).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 10800
      answer:
        has_more: false
        variables: ["DID"]
        results:
        - substitutions:
          - variable: DID
            expression: "'did:example:123456/foo/bar?versionId=1#test'"
      """

  Scenario: Distinguish absent components from empty query and fragment
    This scenario demonstrates that an empty query or fragment is preserved as the empty atom, while absent components remain variables.

    Given the program:
      """ prolog
      empty_components(Query, Fragment, Status) :-
        did_components('did:example:123456?#', did(example, '123456', Path, Query, Fragment)),
        var(Path),
        Status = ok.
      """
    Given the query:
      """ prolog
      consult('/v1/lib/did.pl'),
      empty_components(Query, Fragment, Status).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 7382
      answer:
        has_more: false
        variables: ["Query", "Fragment", "Status"]
        results:
        - substitutions:
          - variable: Query
            expression: "''"
          - variable: Fragment
            expression: "''"
          - variable: Status
            expression: "ok"
      """

  Scenario: Leave missing path, query and fragment unbound
    This scenario demonstrates that optional DID URL components remain unbound when they are absent from the input.

    Given the program:
      """ prolog
      missing_components(Method, MethodSpecificId, Status) :-
        did_components('did:example:123456', did(Method, MethodSpecificId, Path, Query, Fragment)),
        var(Path),
        var(Query),
        var(Fragment),
        Status = ok.
      """
    Given the query:
      """ prolog
      consult('/v1/lib/did.pl'),
      missing_components(Method, MethodSpecificId, Status).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 7176
      answer:
        has_more: false
        variables: ["Method", "MethodSpecificId", "Status"]
        results:
        - substitutions:
          - variable: Method
            expression: "example"
          - variable: MethodSpecificId
            expression: "'123456'"
          - variable: Status
            expression: "ok"
      """

  @great_for_documentation
  Scenario: Error on invalid DID encoding
    This scenario demonstrates the error returned when the DID text does not comply with DID Core syntax.

    Given the query:
      """ prolog
      consult('/v1/lib/did.pl'),
      did_components(foo, Parsed).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4301
      answer:
        has_more: false
        variables: ["Parsed"]
        results:
        - error: "error(domain_error(encoding(did),foo),did_components/2)"
      """

  @great_for_documentation
  Scenario: Error on invalid raw path when reconstructing
    This scenario demonstrates the error returned when a parsed DID term contains a path that is not encoded according to the selected raw representation.

    Given the query:
      """ prolog
      consult('/v1/lib/did.pl'),
      did_components(DID, did(example, '123456', 'path with/space', _, _)).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 6377
      answer:
        has_more: false
        variables: ["DID"]
        results:
        - error: "error(domain_error(encoding(did),path with/space),did_components/2)"
      """

  Scenario: Error on invalid method when reconstructing
    This scenario demonstrates the error returned when the DID method does not satisfy the DID Core grammar.

    Given the query:
      """ prolog
      consult('/v1/lib/did.pl'),
      did_components(DID, did('Example', '123456', _, _, _)).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5073
      answer:
        has_more: false
        variables: ["DID"]
        results:
        - error: "error(domain_error(encoding(did),Example),did_components/2)"
      """

  Scenario: Error on invalid method-specific identifier when reconstructing
    This scenario demonstrates the error returned when the method-specific identifier does not satisfy the DID Core grammar.

    Given the query:
      """ prolog
      consult('/v1/lib/did.pl'),
      did_components(DID, did(example, 'bad id', _, _, _)).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5963
      answer:
        has_more: false
        variables: ["DID"]
        results:
        - error: "error(domain_error(encoding(did),bad id),did_components/2)"
      """

  Scenario: Error on invalid raw query when reconstructing
    This scenario demonstrates the error returned when a parsed DID term contains a query that is not valid raw URI text.

    Given the query:
      """ prolog
      consult('/v1/lib/did.pl'),
      did_components(DID, did(example, '123456', _, 'bad query', _)).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 7207
      answer:
        has_more: false
        variables: ["DID"]
        results:
        - error: "error(domain_error(encoding(did),bad query),did_components/2)"
      """
