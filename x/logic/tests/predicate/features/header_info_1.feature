Feature: header_info/1
  This feature is to test the header_info/1 predicate.

  @great_for_documentation
  Scenario: Retrieve current SDK header info
    This scenario demonstrates how to retrieve the current block header information available to the query.
    The header info contains useful execution context such as the block height, block time, and chain identifier.

    Given the query:
      """ prolog
      consult('/v1/lib/chain.pl'),
      header_info(HeaderInfo).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3980
      answer:
        has_more: false
        variables: ["HeaderInfo"]
        results:
        - substitutions:
          - variable: HeaderInfo
            expression: >-
              header{app_hash:[],chain_id:'axone-testchain-1',hash:[],height:42,time:1712745867}
      """

  @great_for_documentation
  Scenario: Retrieve the block height of the current block.
    This scenario demonstrates how to read the current block height from header_info/1.

    Given a block with the following header:
      """ yaml
      height: 100
      """
    Given the program:
      """ prolog
      :- consult('/v1/lib/chain.pl').

      height(Height) :-
          header_info(Header),
          Height = Header.height.
      """
    Given the query:
      """ prolog
      height(Height).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 100
      gas_used: 3983
      answer:
        has_more: false
        variables: ["Height"]
        results:
        - substitutions:
          - variable: Height
            expression: "100"
      """

  @great_for_documentation
  Scenario: Retrieve the block time of the current block.
    This scenario demonstrates how to read the current block time from header_info/1.

    Given a block with the following header:
      """ yaml
      time: 2024-03-04T11:03:36.000Z
      """
    Given the program:
      """ prolog
      :- consult('/v1/lib/chain.pl').

      time(Time) :-
          header_info(Header),
          Time = Header.time.
      """
    Given the query:
      """ prolog
      time(Time).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3983
      answer:
        has_more: false
        variables: ["Time"]
        results:
        - substitutions:
          - variable: Time
            expression: "1709550216"
      """

  @great_for_documentation
  Scenario: Evaluate a condition based on block time and height
    This scenario demonstrates how to evaluate a condition based on both block time and block height.
    Specifically, it checks whether block time is greater than 1709550216 seconds
    (Monday 4 March 2024 11:03:36 GMT) or block height is greater than 42.

    Given a block with the following header:
      """ yaml
      time: 2024-03-04T11:03:37.000Z
      """
    Given the program:
      """ prolog
      :- consult('/v1/lib/chain.pl').

      evaluate :-
          header_info(Header),
          (Header.time > 1709550216; Header.height > 42),
          !.
      """
    Given the query:
      """ prolog
      evaluate.
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3986
      answer:
        has_more: false
        results:
        - substitutions:
      """
