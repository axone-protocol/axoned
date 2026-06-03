Feature: json_read/2
  This feature validates reading JSON text from a stream into the canonical Prolog JSON representation.

  @great_for_documentation
  Scenario: Read JSON text from a stream
    This scenario demonstrates reading JSON text from a text stream and decoding it into a Prolog term.

    Given the program:
      """ prolog
      :- consult('/v1/lib/json.pl').

      json_read_from_echo(Json, Term) :-
        open('/v1/dev/echo', read_write, Stream, [type(text)]),
        write(Stream, Json),
        json_read(Stream, Term),
        close(Stream).
      """
    Given the query:
      """ prolog
      json_read_from_echo('{"foo":"bar","items":[1,null]}', Term).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 6651
      answer:
        has_more: false
        variables: ["Term"]
        results:
        - substitutions:
          - variable: Term
            expression: "json([foo=bar,items=[1.0,@(null)]])"
      """

  Scenario: Reject malformed JSON read from a stream

    Given the program:
      """ prolog
      :- consult('/v1/lib/json.pl').

      json_read_from_echo(Json, Term) :-
        open('/v1/dev/echo', read_write, Stream, [type(text)]),
        write(Stream, Json),
        json_read(Stream, Term),
        close(Stream).
      """
    Given the query:
      """ prolog
      json_read_from_echo('{&', Term).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5770
      answer:
        has_more: false
        variables: ["Term"]
        results:
        - error: "error(syntax_error(json(malformed_json(1))),json_read/2)"
      """
