Feature: put_char/1
  This feature is to test the put_char/1 predicate.

  @great_for_documentation
  Scenario: Write a single character to user output
  This scenario demonstrates using put_char/1 to write a single character to the current output stream.
  The character appears in the user_output field of the response.

    Given the module configuration:
      """ json
      {
        "limits": {
          "max_user_output_size": 10
        }
      }
      """
    Given the query:
      """ prolog
      put_char('b').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3986
      answer:
        has_more: false
        variables:
        results:
        - substitutions:
      user_output: "b"
      """

  @great_for_documentation
  Scenario: Write multiple characters to user output
  This scenario demonstrates chaining multiple put_char/1 calls to write several characters.
  Each character is appended to the user output stream.

    Given the module configuration:
      """ json
      {
        "limits": {
          "max_user_output_size": 10
        }
      }
      """
    Given the query:
      """ prolog
      put_char('a'), put_char('b'), put_char('c').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4042
      answer:
        has_more: false
        variables:
        results:
        - substitutions:
      user_output: "abc"
      """

  @great_for_documentation
  Scenario: Write characters with user output size limit
  This scenario shows how the user output is truncated when it exceeds the configured max_user_output_size limit.
  The limit is measured in bytes, so only the last characters that fit within the limit are kept.

    Given the module configuration:
      """ json
      {
        "limits": {
          "max_user_output_size": 3
        }
      }
      """
    Given the query:
      """ prolog
      put_char('h'), put_char('e'), put_char('l'), put_char('l'), put_char('o').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4098
      answer:
        has_more: false
        variables:
        results:
        - substitutions:
      user_output: "llo"
      """

  @great_for_documentation
  Scenario: Write UTF-8 character
  This scenario illustrates writing UTF-8 characters using put_char/1.
  Multi-byte characters like emojis occupy more space in the buffer.

    Given the module configuration:
      """ json
      {
        "limits": {
          "max_user_output_size": 10
        }
      }
      """
    Given the program:
      """ prolog
      log_message([]).
      log_message([H|T]) :-
          put_char(H),
          log_message(T).
      """
    Given the query:
      """ prolog
      log_message("😀").
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4093
      answer:
        has_more: false
        variables:
        results:
        - substitutions:
      user_output: "😀"
      """

  Scenario: Write newline character
  This scenario demonstrates writing special characters like newline using put_char/1.

    Given the module configuration:
      """ json
      {
        "limits": {
          "max_user_output_size": 20
        }
      }
      """
    Given the query:
      """ prolog
      put_char('a'), put_char('\n'), put_char('b').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4043
      answer:
        has_more: false
        variables:
        results:
        - substitutions:
      user_output: |
        a
        b
      """
