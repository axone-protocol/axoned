Feature: json_write/2
  This feature validates writing canonical Prolog JSON terms as JSON text.

  @great_for_documentation
  Scenario: Write a canonical Prolog JSON term to a stream
    This scenario demonstrates writing a Prolog JSON term to the current output stream.

    Given the module configuration:
      """ json
      {
        "limits": {
          "max_user_output_size": 50
        }
      }
      """
    Given the program:
      """ prolog
      :- consult('/v1/lib/json.pl').

      json_write_to_output(Term) :-
        current_output(Stream),
        json_write(Stream, Term).
      """
    Given the query:
      """ prolog
      json_write_to_output(json([foo=bar,ok= @(true)])).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 11668
      answer:
        has_more: false
        variables:
        results:
        - substitutions:
      user_output: |
        {"foo":"bar","ok":true}
      """

  Scenario: Reject invalid Prolog JSON terms

    Given the module configuration:
      """ json
      {
        "limits": {
          "max_user_output_size": 50
        }
      }
      """
    Given the program:
      """ prolog
      :- consult('/v1/lib/json.pl').

      json_write_to_output(Term) :-
        current_output(Stream),
        json_write(Stream, Term).
      """
    Given the query:
      """ prolog
      json_write_to_output(foo([a=b])).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 7991
      answer:
        has_more: false
        variables:
        results:
        - error: "error(type_error(json,foo([=(a,b)])),json_write/2)"
      user_output: ""
      """
