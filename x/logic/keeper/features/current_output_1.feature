Feature: current_output/1
  This feature is to test the current_output/1 predicate.

  @great_for_documentation
  Scenario: Write a char to the current output
  This scenario demonstrates how to write a character to the current output, and get the content in the response of the
  request.

    Given the module configuration:
      """ json
      {
        "limits": {
          "max_user_output_size": "5"
        }
      }
      """
    Given the program:
      """ prolog
      write_char_to_user_output(C) :-
          current_output(UserStream), % get the current output stream
          put_char(UserStream, C).    % write the char to the user stream
      """
    Given the query:
      """ prolog
      write_char_to_user_output(x).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4241
      answer:
        has_more: false
        variables:
        results:
        - substitutions:
      user_output: |
        x
      """

  @great_for_documentation
  Scenario: Write characters to the current output (without limit)
  This scenario demonstrates how to write some characters to the current output, and get the content in the response of the
  request. This is helpful for debugging purposes.

    Given the module configuration:
      """ json
      {
        "limits": {
          "max_user_output_size": "15"
        }
      }
      """
    Given the program:
      """ prolog
      log_message(Message) :-
          current_output(UserStream), % get the current output stream
          write(UserStream, Message), % write the message to the user stream
          put_char(UserStream, '\n').
      """
    Given the query:
      """ prolog
      log_message('Hello world!').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4276
      answer:
        has_more: false
        variables:
        results:
        - substitutions:
      user_output: |
        Hello world!

      """

  @great_for_documentation
  Scenario: Write characters to the current output (with limit)
  This scenario demonstrates the process of writing characters to the current user output, with a limit configured
  in the logic module. So if the message is longer than this limit, the output will be truncated.

    Given the module configuration:
      """ json
      {
        "limits": {
          "max_user_output_size": "5"
        }
      }
      """
    Given the program:
      """ prolog
      log_message(Message) :-
          current_output(UserStream), % get the current output stream
          write(UserStream, Message). % write the message to the user stream
      """
    Given the query:
      """ prolog
      log_message('Hello world!').
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4242
      answer:
        has_more: false
        variables:
        results:
        - substitutions:
      user_output: |
        orld!
      """

  @great_for_documentation
  Scenario: Write UTF-8 character to the current output (with limit)
  This scenario illustrates the impact of UTF-8 characters on output limits measured in bytes, not character count.
  Characters such as emojis require more space; for example, the wizard emoji (🧙) occupies 4 bytes, effectively counting
  as four units. As a result, the limit is reached more quickly with these characters, which means that the number of
  characters in the user output is less than expected.

    Given the module configuration:
      """ json
      {
        "limits": {
          "max_user_output_size": "5"
        }
      }
      """
    Given the program:
      """ prolog
      log_message([]).
      log_message([H|T]) :-
          current_output(UserStream),
          put_char(UserStream, H),
          log_message(T).
      """
    Given the query:
      """ prolog
      log_message("Hello 🧙!").
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4263
      answer:
        has_more: false
        variables:
        results:
        - substitutions:
      user_output: "🧙!"
      """

  Scenario: Write strings to the current output (no limit configured)
  This scenario demonstrates that if no limit is configured in the logic module, the user can write as much as they want.
  This case should not be used in production, as it can lead to performance issues.

    Given the program:
      """ prolog
      log_message([]).
      log_message([H|T]) :-
          current_output(UserStream),
          put_char(UserStream, H),
          log_message(T).
      """
    Given the query:
      """ prolog
      log_message("Prolog's logic weaves through the fabric of the chain,\nGovernance and rules, in its domain reign."),
      log_message("\n"),
      log_message("Knowledge blooms in the heart of the AXONE lore,\nUnlocking a world of possibilities to explore.").
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 4721
      answer:
        has_more: false
        variables:
        results:
        - substitutions:
      user_output: |
        Prolog's logic weaves through the fabric of the chain,
        Governance and rules, in its domain reign.
        Knowledge blooms in the heart of the AXONE lore,
        Unlocking a world of possibilities to explore.
      """
