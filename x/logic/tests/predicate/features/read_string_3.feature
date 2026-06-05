Feature: read_string/3
  This feature is to test the read_string/3 predicate.

  @great_for_documentation
  Scenario: Read a text stream into an atom and byte length
  This scenario demonstrates reading all text from a stream while counting UTF-8 bytes.

    Given the program:
      """ prolog
      read_from_echo(Text, Length, String) :-
        open('/v1/dev/echo', read_write, Stream, [type(text)]),
        write(Stream, Text),
        read_string(Stream, Length, String),
        close(Stream).
      """
    Given the query:
      """ prolog
      read_from_echo('aé', Length, String).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5505
      answer:
        has_more: false
        variables: ["Length", "String"]
        results:
        - substitutions:
          - variable: Length
            expression: 3
          - variable: String
            expression: "aé"
      """

  Scenario: Fail when a fixed length is longer than the stream

    Given the program:
      """ prolog
      read_from_echo(Text, Length, String) :-
        open('/v1/dev/echo', read_write, Stream, [type(text)]),
        write(Stream, Text),
        read_string(Stream, Length, String),
        close(Stream).
      """
    Given the query:
      """ prolog
      read_from_echo(abc, 5, String).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5980
      answer:
        has_more: false
        variables: ["String"]
        results:
      """
