Feature: json_prolog/2
  This feature validates JSON conversion between text and the canonical Prolog JSON representation.

  @great_for_documentation
  Scenario: Decode JSON text into a canonical Prolog term
    This scenario demonstrates how JSON objects, strings, and booleans are represented in Prolog.

    Given the query:
      """ prolog
      consult('/v1/lib/json.pl'),
      json_prolog('{"foo":"bar","ok":true}', Term).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 6216
      answer:
        has_more: false
        variables: ["Term"]
        results:
        - substitutions:
          - variable: Term
            expression: "json([foo=bar,ok= @(true)])"
      """

  @great_for_documentation
  Scenario: Encode a canonical Prolog term as JSON text
    This scenario demonstrates how a canonical Prolog JSON object is encoded as compact JSON text.

    Given the query:
      """ prolog
      consult('/v1/lib/json.pl'),
      json_prolog(Json, json([foo=bar,ok= @(true)])).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 11049
      answer:
        has_more: false
        variables: ["Json"]
        results:
        - substitutions:
          - variable: Json
            expression: "'{\"foo\":\"bar\",\"ok\":true}'"
      """

  Scenario: Reject malformed JSON text

    Given the query:
      """ prolog
      consult('/v1/lib/json.pl'),
      json_prolog('{&', Term).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5560
      answer:
        has_more: false
        variables: ["Term"]
        results:
        - error: "error(syntax_error(json(malformed_json(1))),json_prolog/2)"
      """

  Scenario: Reject invalid Prolog JSON terms

    Given the query:
      """ prolog
      consult('/v1/lib/json.pl'),
      json_prolog(Json, foo([a=b])).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 8034
      answer:
        has_more: false
        variables: ["Json"]
        results:
        - error: "error(type_error(json,foo([=(a,b)])),json_prolog/2)"
      """

  Scenario: Require one instantiated side

    Given the query:
      """ prolog
      consult('/v1/lib/json.pl'),
      json_prolog(Json, Term).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4191
      answer:
        has_more: false
        variables: ["Json", "Term"]
        results:
        - error: "error(instantiation_error,json_prolog/2)"
      """
