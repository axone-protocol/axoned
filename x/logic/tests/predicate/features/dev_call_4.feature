Feature: dev_call/4
  This feature validates the transactional device helper predicate, which provides
  a high-level interface for interacting with half-duplex transactional devices.

  A device implements a half-duplex protocol with three phases:
  1. Request phase: write operations build up the request payload
  2. Commit phase: the first read operation commits the request
  3. Response phase: subsequent reads stream the response until EOF

  Once committed (after first read), write operations are rejected.

  Background:
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

  @great_for_documentation
  Scenario: Echo roundtrip using dev_call/4 with meta goals
    This scenario demonstrates the typical successful usage of dev_call/4:
    - The write goal sends request bytes to the device
    - The read goal commits the request and reads the response
    - The device echoes back the exact bytes sent

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

  @great_for_documentation
  Scenario: Reading without writing fails with invalid_request
    This scenario illustrates the half-duplex protocol requirement:
    the device expects at least one write before the first read commits.

    Given the program:
      """ prolog
      :- consult('/v1/lib/dev.pl').

      read_without_write(Result) :-
        dev_call('/v1/dev/echo', binary, no_write, read_bytes(Result)).

      no_write(_).

      read_bytes(Stream, Bytes) :-
        dev_read_bytes(Stream, Bytes).
      """
    Given the query:
      """ prolog
      catch(
        read_without_write(_),
        error(system_error, _),
        true
      ).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5021
      answer:
        has_more: false
        results:
        - {}
      """

  @great_for_documentation
  Scenario: Writing after reading fails with permission error
    This scenario demonstrates that once the first read commits the request,
    the device transitions to read-only mode and rejects further writes.

    Given the program:
      """ prolog
      :- consult('/v1/lib/dev.pl').

      write_after_read(Echoed) :-
        dev_call('/v1/dev/echo', binary, write_then_read(Echoed), no_read).

      no_read(_).

      write_then_read(Stream, Echoed) :-
        dev_write_bytes(Stream, [1,2,3]),
        dev_read_bytes(Stream, Echoed),
        dev_write_bytes(Stream, [4,5,6]).
      """
    Given the query:
      """ prolog
      catch(
        write_after_read(_),
        error(system_error, _),
        true
      ).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 6415
      answer:
        has_more: false
        results:
        - {}
      """

  @great_for_documentation
  Scenario: Multiple reads stream the response progressively
    This scenario shows that after commit, multiple read operations
    can progressively consume the response stream until EOF.

    Given the program:
      """ prolog
      :- consult('/v1/lib/dev.pl').

      echo_partial(Result) :-
        dev_call('/v1/dev/echo', binary, write_all, read_partial(Result)).

      write_all(Stream) :-
        dev_write_bytes(Stream, [65,66,67,68]).

      read_partial(Stream, Result) :-
        get_byte(Stream, B1),
        get_byte(Stream, B2),
        Result = [B1, B2].
      """
    Given the query:
      """ prolog
      echo_partial(Result).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5803
      answer:
        has_more: false
        variables: ["Result"]
        results:
        - substitutions:
          - variable: Result
            expression: "[65,66]"
      """
