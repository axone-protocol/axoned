Feature: uri_encoded/3
  This feature is to test the uri_encoded/3 predicate.

  @great_for_documentation
  Scenario: Decode a raw path atom
    This scenario demonstrates how to decode a raw URI path into plain text.

    Given the query:
      """ prolog
      consult('/v1/lib/uri.pl'),
      uri_encoded(path, Decoded, foo).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4667
      answer:
        has_more: false
        variables: ["Decoded"]
        results:
        - substitutions:
          - variable: Decoded
            expression: "foo"
      """

  @great_for_documentation
  Scenario: Encode a query value with a space
    This scenario demonstrates how to percent-encode a query value.

    Given the query:
      """ prolog
      consult('/v1/lib/uri.pl'),
      uri_encoded(query_value, 'foo bar', Encoded).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5702
      answer:
        has_more: false
        variables: ["Encoded"]
        results:
        - substitutions:
          - variable: Encoded
            expression: "'foo%20bar'"
      """

  Scenario: Encode component-specific reserved characters for a query value
    This scenario demonstrates the reserved-character policy specific to the query_value component.

    Given the query:
      """ prolog
      consult('/v1/lib/uri.pl'),
      uri_encoded(query_value, '&+/:;=?[]^{}', Encoded).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 8367
      answer:
        has_more: false
        variables: ["Encoded"]
        results:
        - substitutions:
          - variable: Encoded
            expression: "'%26%2B/%3A%3B%3D?%5B%5D%5E%7B%7D'"
      """

  Scenario: Encode component-specific reserved characters for a path
    This scenario demonstrates the reserved-character policy specific to the path component.

    Given the query:
      """ prolog
      consult('/v1/lib/uri.pl'),
      uri_encoded(path, ':/?[]^{}', Encoded).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 7066
      answer:
        has_more: false
        variables: ["Encoded"]
        results:
        - substitutions:
          - variable: Encoded
            expression: "'%3A/%3F%5B%5D%5E%7B%7D'"
      """

  Scenario: Encode component-specific reserved characters for a segment
    This scenario demonstrates the reserved-character policy specific to the segment component.

    Given the query:
      """ prolog
      consult('/v1/lib/uri.pl'),
      uri_encoded(segment, ':/?[]^{}', Encoded).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 7215
      answer:
        has_more: false
        variables: ["Encoded"]
        results:
        - substitutions:
          - variable: Encoded
            expression: "'%3A%2F%3F%5B%5D%5E%7B%7D'"
      """

  Scenario: Encode component-specific reserved characters for a fragment
    This scenario demonstrates the reserved-character policy specific to the fragment component.

    Given the query:
      """ prolog
      consult('/v1/lib/uri.pl'),
      uri_encoded(fragment, '<>?[]^{}', Encoded).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 7054
      answer:
        has_more: false
        variables: ["Encoded"]
        results:
        - substitutions:
          - variable: Encoded
            expression: "'%3C%3E?%5B%5D%5E%7B%7D'"
      """

  Scenario: Decode component-specific reserved characters for a query value
    This scenario demonstrates that percent-decoding uses generic URI unescaping and preserves raw plus characters.

    Given the query:
      """ prolog
      consult('/v1/lib/uri.pl'),
      uri_encoded(query_value, Decoded, '%26%2B/%3A%3B%3D?%5B%5D%5E%7B%7D').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 6172
      answer:
        has_more: false
        variables: ["Decoded"]
        results:
        - substitutions:
          - variable: Decoded
            expression: "'&+/:;=?[]^{}'"
      """

  Scenario: Error on an invalid URI component
    This scenario demonstrates the error returned when the requested URI component is unknown.

    Given the query:
      """ prolog
      consult('/v1/lib/uri.pl'),
      uri_encoded(hey, foo, Decoded).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4106
      answer:
        has_more: false
        variables: ["Decoded"]
        results:
        - error: "error(type_error(uri_component,hey),uri_encoded/3)"
      """

  Scenario: Error on an unbound URI component
    This scenario demonstrates the error returned when the URI component itself is not instantiated.

    Given the query:
      """ prolog
      consult('/v1/lib/uri.pl'),
      uri_encoded(Component, 'foo bar', 'bar%20foo').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4100
      answer:
        has_more: false
        variables: ["Component"]
        results:
        - error: "error(instantiation_error,uri_encoded/3)"
      """

  Scenario: Error on an invalid text value in encode mode
    This scenario demonstrates the error returned when the plain-text value is not a text term.

    Given the query:
      """ prolog
      consult('/v1/lib/uri.pl'),
      uri_encoded(path, compound(2), 'bar%20foo').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4644
      answer:
        has_more: false
        results:
        - error: "error(type_error(text,compound(2)),uri_encoded/3)"
      """

  Scenario: Error on an invalid text value in decode mode
    This scenario demonstrates the error returned when the encoded value is not a text term.

    Given the query:
      """ prolog
      consult('/v1/lib/uri.pl'),
      uri_encoded(path, Decoded, compound(2)).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4720
      answer:
        has_more: false
        variables: ["Decoded"]
        results:
        - error: "error(type_error(text,compound(2)),uri_encoded/3)"
      """

  Scenario: Fail on mismatching fully ground values
    This scenario demonstrates that a fully ground mismatch simply fails.

    Given the query:
      """ prolog
      consult('/v1/lib/uri.pl'),
      uri_encoded(fragment, 'foo bar', 'bar%20foo').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5719
      answer:
        has_more: false
      """

  Scenario: Fail when the encoded target is a non-unifiable non-text term in encode mode
    This scenario demonstrates that encode mode computes the encoded atom and then simply fails on unification.

    Given the query:
      """ prolog
      consult('/v1/lib/uri.pl'),
      uri_encoded(path, foo, compound(2)).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4880
      answer:
        has_more: false
      """

  Scenario: Error on an invalid percent escape
    This scenario demonstrates the error returned when the encoded input contains an invalid percent escape sequence.

    Given the query:
      """ prolog
      consult('/v1/lib/uri.pl'),
      uri_encoded(path, Decoded, 'bar%%3foo').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5166
      answer:
        has_more: false
        variables: ["Decoded"]
        results:
        - error: "error(domain_error(encoding(uri),bar%%3foo),uri_encoded/3,[i,n,v,a,l,i,d, ,U,R,L, ,e,s,c,a,p,e, ,\",%,%,3,\"])"
      """
