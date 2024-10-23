Feature: block_height/1
  This feature is to test the block_height/2 predicate.

  @great_for_documentation
  Scenario: Retrieve the block height of the current block.
    This scenario demonstrates how to retrieve the block height of the current block.

    Given a block with the following header:
      | Height | 100   |

    Given the query:
      """ prolog
      block_height(Height).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 100
      gas_used: 3975
      answer:
        has_more: false
        variables: ["Height"]
        results:
        - substitutions:
          - variable: Height
            expression: "100"
      """

  @great_for_documentation
  Scenario: Check that the block height is greater than a certain value.
    This scenario demonstrates how to check that the block height is greater than 100. This predicate is useful for
    governance which requires a certain block height to be reached before a certain action is taken.

    Given a block with the following header:
      | Height | 101   |

    Given the query:
      """ prolog
      block_height(Height),
      Height > 100.
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 101
      gas_used: 3976
      answer:
        has_more: false
        variables: ["Height"]
        results:
        - substitutions:
          - variable: Height
            expression: "101"
      """
