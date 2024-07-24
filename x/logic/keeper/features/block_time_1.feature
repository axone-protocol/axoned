Feature: block_time/1
  This feature is to test the block_time/2 predicate.

  @great_for_documentation
  Scenario: Retrieve the block time of the current block.
    This scenario demonstrates how to retrieve the block time of the current block.

    Given a block with the following header:
      | Time | 1709550216   |

    Given the query:
      """ prolog
      block_time(Time).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3876
      answer:
        has_more: false
        variables: ["Time"]
        results:
        - substitutions:
          - variable: Time
            expression: "1709550216"
      """

  @great_for_documentation
  Scenario: Check that the block time is greater than a certain time.
    This scenario demonstrates how to check that the block time is greater than 1709550216 seconds (Monday 4 March 2024 11:03:36 GMT)
    using the `block_time/1` predicate. This predicate is useful for governance which requires a certain block time to be
    reached before a certain action is taken.

    Given a block with the following header:
      | Time | 1709550217   |

    Given the query:
      """ prolog
      block_time(Time),
      Time > 1709550216.
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3877
      answer:
        has_more: false
        variables: ["Time"]
        results:
        - substitutions:
          - variable: Time
            expression: "1709550217"
      """
