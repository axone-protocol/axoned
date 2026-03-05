Feature: dev_call/4
  This feature validates the transactional device helper and meta-goal wiring.

  Scenario: Echo roundtrip using dev_call/4 with meta goals
    Given the program:
      """ prolog
      :- consult('/v1/lib/dev.pl').

      echo(Bytes, Echoed) :-
        dev_call('/v1/dev/echo', binary, write_bytes(Bytes), read_bytes(Echoed)).

      write_bytes(Stream, Bytes) :-
        dev_write_bytes(Stream, Bytes).

      read_bytes(Stream, Bytes) :-
        dev_read_bytes(Stream, Bytes).
      """
    Given the query:
      """ prolog
      echo([0,1,2,255], Echoed).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 6270
      answer:
        has_more: false
        variables: ["Echoed"]
        results:
        - substitutions:
          - variable: Echoed
            expression: "[0,1,2,255]"
      """
