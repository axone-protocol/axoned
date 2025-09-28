Feature: open/3
  This feature is to test the open/3 predicate.

  @great_for_documentation
  Scenario: Open a resource for reading
  This scenario showcases the procedure for accessing a resource stored within a CosmWasm smart contract for reading
  purposes and obtaining the stream's properties.

  See the `open/4` predicate for a more detailed example.

    Given the CosmWasm smart contract "axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk" and the behavior:
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
      open(
        'cosmwasm:storage:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D&base64Decode=false',
        read,
        _
      ).
      """
    When the query is run
    Then the answer we get is:
      """ yaml
      height: 42
      gas_used: 3976
      answer:
        has_more: false
        variables:
        results:
        - substitutions:
      """
