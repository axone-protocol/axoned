Feature: wasm_query/3
  This feature is to test the wasm_query/3 predicate.

  @great_for_documentation
  Scenario: Query a smart contract and read raw response bytes
    This scenario demonstrates how to send a smart-query payload as bytes and get the raw response bytes back.

    Given the CosmWasm smart contract "axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk" and the behavior:
      """ yaml
      message: |
        {"ping":"pong"}
      response: |
        {"ok":true}
      """
    Given the program:
      """ prolog
      :- consult('/v1/lib/wasm.pl').
      """
    Given the query:
      """ prolog
      wasm_query(
        'axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk',
        [123,34,112,105,110,103,34,58,34,112,111,110,103,34,125],
        ResponseBytes
      ).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4478
      answer:
        has_more: false
        variables: ["ResponseBytes"]
        results:
        - substitutions:
          - variable: ResponseBytes
            expression: "[123,34,111,107,34,58,116,114,117,101,125]"
      """

  @great_for_documentation
  Scenario: Reject an invalid contract address
    This scenario demonstrates how wasm_query/3 validates the contract address format before calling the chain.

    Given the program:
      """ prolog
      :- consult('/v1/lib/wasm.pl').
      """
    Given the query:
      """ prolog
      wasm_query('invalid-address', [123,125], ResponseBytes).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4000
      answer:
        has_more: false
        variables: ["ResponseBytes"]
        results:
        - error: "error(domain_error(encoding(bech32),invalid-address),wasm_query/3)"
      """

  @great_for_documentation
  Scenario: Reject a non-byte request payload
    This scenario demonstrates how wasm_query/3 rejects payload lists containing values outside the byte range [0,255].

    Given the program:
      """ prolog
      :- consult('/v1/lib/wasm.pl').
      """
    Given the query:
      """ prolog
      wasm_query(
        'axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk',
        [256],
        ResponseBytes
      ).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4065
      answer:
        has_more: false
        variables: ["ResponseBytes"]
        results:
        - error: "error(type_error(byte,256),must_be/2)"
      """

  @great_for_documentation
  Scenario: Surface contract query execution failures
    This scenario demonstrates that a contract query failure is returned as a system_error.

    Given the CosmWasm smart contract "axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk" and the behavior:
      """ yaml
      message: |
        {"ping":"pong"}
      error: wasm contract execution failed
      """
    Given the program:
      """ prolog
      :- consult('/v1/lib/wasm.pl').
      """
    Given the query:
      """ prolog
      wasm_query(
        'axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk',
        [123,34,112,105,110,103,34,58,34,112,111,110,103,34,125],
        ResponseBytes
      ).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4397
      answer:
        has_more: false
        variables: ["ResponseBytes"]
        results:
        - error: "error(system_error,read /v1/dev/wasm/axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk/query: wasm_query_failed)"
      """
