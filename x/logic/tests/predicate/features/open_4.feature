Feature: open/4
  This feature is to test the open/4 predicate.

  @great_for_documentation
  Scenario: Open a resource for reading
  This scenario demonstrates how to build a VFS path and open it in read mode.

    Given the program:
      """ prolog
      resource_path(Path) :-
        atomic_list_concat(['/v1', 'sys', 'header', 'height'], '/', Path).
      """
    Given the query:
      """ prolog
      resource_path(Path),
      open(Path, read, _, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 6864
      answer:
        has_more: false
        variables: ["Path"]
        results:
        - substitutions:
          - variable: Path
            expression: "'/v1/sys/header/height'"
      """

  @great_for_documentation
  Scenario: Open an existing resource and read its Prolog term
  This scenario demonstrates how to open a text resource, read one term, and close the stream.

    Given the program:
      """ prolog
      read_height(Path, Height) :-
        open(Path, read, Stream, [type(text)]),
        read_term(Stream, Height, []),
        close(Stream).
      """
    Given the query:
      """ prolog
      read_height('/v1/sys/header/height', Height).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4100
      answer:
        has_more: false
        variables: ["Height"]
        results:
        - substitutions:
          - variable: Height
            expression: "42"
      """

  @great_for_documentation
  Scenario: Open a wasm query endpoint in read_write mode
  This scenario demonstrates how to write request bytes, then read response bytes from a transactional endpoint.

    Given the CosmWasm smart contract "axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk" and the behavior:
      """ yaml
      message: |
        {}
      response: |
        {"ok":true}
      """
    Given the program:
      """ prolog
      read_all_bytes(Stream, Bytes) :-
        get_byte(Stream, Byte),
        ( Byte =:= -1 ->
            Bytes = []
        ; Bytes = [Byte | Rest],
          read_all_bytes(Stream, Rest)
        ).

      wasm_roundtrip(Address, ResponseBytes) :-
        atom_concat('/v1/dev/wasm/', Address, Prefix),
        atom_concat(Prefix, '/query', Path),
        open(Path, read_write, Stream, [type(binary)]),
        put_byte(Stream, 123),
        put_byte(Stream, 125),
        read_all_bytes(Stream, ResponseBytes),
        close(Stream).
      """
    Given the query:
      """ prolog
      wasm_roundtrip('axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk', ResponseBytes).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 5641
      answer:
        has_more: false
        variables: ["ResponseBytes"]
        results:
        - substitutions:
          - variable: ResponseBytes
            expression: "[123,34,111,107,34,58,116,114,117,101,125]"
      """

  @great_for_documentation
  Scenario: Try to open a non-existing resource
    This scenario demonstrates the system's response to trying to open a non-existing resource.

    Given the query:
      """ prolog
      open('foo:bar', read, Stream, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3935
      answer:
        has_more: false
        variables: ["Stream"]
        results:
        - error: "error(existence_error(source_sink,foo:bar),open/4)"
      """

  @great_for_documentation
  Scenario: Try to open a read-only resource for writing
  This scenario demonstrates the system's response to opening a snapshot path in write mode.

    Given the query:
      """ prolog
      open('/v1/sys/header/height', write, Stream, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3950
      answer:
        has_more: false
        variables: ["Stream"]
        results:
        - error: "error(permission_error(open,source_sink,/v1/sys/header/height),open/4)"
      """

  @great_for_documentation
  Scenario: Try to open a read-only resource for appending
  This scenario demonstrates the system's response to opening a snapshot path in append mode.

    Given the query:
      """ prolog
      open('/v1/sys/header/height', append, Stream, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3951
      answer:
        has_more: false
        variables: ["Stream"]
        results:
        - error: "error(permission_error(open,source_sink,/v1/sys/header/height),open/4)"
      """

  @great_for_documentation
  Scenario: Pass incorrect options to open/4
  This scenario demonstrates the system's response to opening a resource with incorrect options.

    Given the query:
      """ prolog
      open('/v1/sys/header/height', read, Stream, [non_existing_option]).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3971
      answer:
        has_more: false
        variables: ["Stream"]
        results:
        - error: "error(domain_error(stream_option,non_existing_option),open/4)"
      """

  Scenario: Open a resource with incorrect mode (1)
  This scenario demonstrates the system's response to opening a resource with an incorrect mode.

    Given the query:
      """ prolog
      open('/v1/sys/header/height', incorrect_mode, Stream, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3959
      answer:
        has_more: false
        variables: ["Stream"]
        results:
        - error: "error(domain_error(io_mode,incorrect_mode),open/4)"
      """

  Scenario: Open a resource with incorrect mode (2)
  This scenario demonstrates the system's response to opening a resource with an incorrect mode.

    Given the query:
      """ prolog
      open('/v1/sys/header/height', 666, Stream, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3948
      answer:
        has_more: false
        variables: ["Stream"]
        results:
        - error: "error(type_error(atom,666),open/4)"
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
      height: 42
      gas_used: 3934
      answer:
        has_more: false
        variables: ["Resource", "Stream"]
        results:
        - error: "error(instantiation_error,open/4)"
      """

  Scenario: Insufficient instantiation error (2)
  This scenario demonstrates the system's response to calling open/4 with insufficiently instantiated arguments.

    Given the query:
      """ prolog
      open('/v1/sys/header/height', Mode, Stream, []).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3949
      answer:
        has_more: false
        variables: ["Mode", "Stream"]
        results:
        - error: "error(instantiation_error,open/4)"
      """
