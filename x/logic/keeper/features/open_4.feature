Feature: open/4
  This feature is to test the open/4 predicate.

  @great_for_documentation
  Scenario: Open a resource for reading
  This scenario showcases the procedure for accessing a resource stored within a CosmWasm smart contract for reading
  purposes and obtaining the stream's properties.

  Assuming the existence of a CosmWasm smart contract configured to store resources, we construct a URI to specifically
  identify the smart contract and pinpoint the resource we aim to retrieve via a query message.

    Given the CosmWasm smart contract "okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht" and the behavior:
      """ yaml
      message: |
        {
          "object_data": {
            "id": "4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05"
          }
        }
      response: |
        Hello, World!
      """
    Given the query:
      """ prolog
      open('cosmwasm:storage:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D&base64Decode=false', read, _, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      has_more: false
      variables:
      results:
      - substitutions:
      """

  @great_for_documentation
  Scenario: Open an existing resource and read its content
  This scenario shows a more complex example of how to open an existing resource stored in a CosmWasm smart contract
  and read its content.

  The resource is opened for reading, and the content is read into a list of characters. Finally, the stream is closed.

    Given the CosmWasm smart contract "okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht" and the behavior:
      """ yaml
      message: |
        {
          "object_data": {
            "id": "4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05"
          }
        }
      response: |
        Hello, World!
      """
    Given the program:
      """ prolog
      read_resource(Resource, Chars) :-
        open('cosmwasm:storage:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D&base64Decode=false', read, Stream, []),
        get_char(Stream, Char),
        read_resource(Stream, Char, Chars),
        close(Stream).

      read_resource(Stream, Char, Chars) :-
        (   Char == end_of_file ->
            Chars = []
        ;   Chars = [Char| Rest],
            get_char(Stream, Next),
            read_resource(Stream, Next, Rest)
        ).
      """
    Given the query:
      """ prolog
      read_resource('hello-world.txt', Chars).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      has_more: false
      variables: ["Chars"]
      results:
      - substitutions:
        - variable: Chars
          expression: "['H',e,l,l,o,',',' ','W',o,r,l,d,!]"
      """

  @great_for_documentation
  Scenario: Open an existing resource and read its content (base64-encoded)
  This scenario is a variation of the previous one. The difference is that the smart contract returns a base64-encoded
  response. For this reason, we set the `base64Decode` parameter to `true` in the query (the default value is `false`).

    Given the CosmWasm smart contract "okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht" and the behavior:
      """ yaml
      message: |
        {
          "object_data": {
            "id": "4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05"
          }
        }
      response: |
        "SGVsbG8sIFdvcmxkIQ=="
      """
    Given the program:
      """ prolog
      read_resource(Resource, Chars) :-
        open('cosmwasm:storage:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D&base64Decode=true', read, Stream, []),
        get_char(Stream, Char),
        read_resource(Stream, Char, Chars),
        close(Stream).

      read_resource(Stream, Char, Chars) :-
        (   Char == end_of_file ->
            Chars = []
        ;   Chars = [Char| Rest],
            get_char(Stream, Next),
            read_resource(Stream, Next, Rest)
        ).
      """
    Given the query:
      """ prolog
      read_resource('hello-world.txt', Chars).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      has_more: false
      variables: ["Chars"]
      results:
      - substitutions:
        - variable: Chars
          expression: "['H',e,l,l,o,',',' ','W',o,r,l,d,!]"
      """

  @great_for_documentation
  Scenario: Try to open a non-existing resource
    This scenario demonstrates the system's response to trying to open a non-existing resource.

    Given the query:
      """ prolog
      open('cosmwasm:storage:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=foo', read, Stream, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      has_more: false
      variables: ["Stream"]
      results:
      - error: "error(existence_error(source_sink,cosmwasm:storage:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=foo),open/4)"
      """

  @great_for_documentation
  Scenario: Try to open a resource for writing
  This scenario demonstrates the system's response to opening a resource for writing, but the resource does not allow
  writing. This is the case for resources hosted in smart contracts which are read-only.

    Given the query:
      """ prolog
      open('cosmwasm:storage:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=foo', write, Stream, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      has_more: false
      variables: ["Stream"]
      results:
      - error: "error(permission_error(input,stream,cosmwasm:storage:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=foo),open/4)"
      """

  @great_for_documentation
  Scenario: Try to open a resource for appending
  This scenario demonstrates the system's response to opening a resource for appending, but the resource does not allow
  appending. This is the case for resources hosted in smart contracts which are read-only.

    Given the query:
      """ prolog
      open('cosmwasm:storage:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=foo', write, Stream, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      has_more: false
      variables: ["Stream"]
      results:
      - error: "error(permission_error(input,stream,cosmwasm:storage:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=foo),open/4)"
      """


  @great_for_documentation
  Scenario: Pass incorrect options to open/4
  This scenario demonstrates the system's response to opening a resource with incorrect options.

    Given the query:
      """ prolog
      open('cosmwasm:storage:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=foo', read, Stream, [non_existing_option]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      has_more: false
      variables: ["Stream"]
      results:
      - error: "error(domain_error(empty_list,[non_existing_option]),open/4)"
      """


  Scenario: Open a resource with incorrect mode (1)
  This scenario demonstrates the system's response to opening a resource with an incorrect mode.

    Given the query:
      """ prolog
      open('cosmwasm:storage:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=foo', incorrect_mode, Stream, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      has_more: false
      variables: ["Stream"]
      results:
      - error: "error(type_error(io_mode,incorrect_mode),open/4)"
      """

  Scenario: Open a resource with incorrect mode (2)
  This scenario demonstrates the system's response to opening a resource with an incorrect mode.

    Given the query:
      """ prolog
      open('cosmwasm:storage:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=foo', 666, Stream, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      has_more: false
      variables: ["Stream"]
      results:
      - error: "error(type_error(io_mode,666),open/4)"
      """

  Scenario: Insufficient instantiation error (1)
  This scenario demonstrates the system's response to calling open/4 with insufficiently instantiated arguments.

    Given the query:
      """ prolog
      open(Resource, read, Stream, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      has_more: false
      variables: ["Resource", "Stream"]
      results:
      - error: "error(instantiation_error,open/4)"
      """

  Scenario: Insufficient instantiation error (2)
  This scenario demonstrates the system's response to calling open/4 with insufficiently instantiated arguments.

    Given the query:
      """ prolog
      open('cosmwasm:storage:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=foo', Mode, Stream, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      has_more: false
      variables: ["Mode", "Stream"]
      results:
      - error: "error(instantiation_error,open/4)"
      """
