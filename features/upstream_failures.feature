Feature: Upstream failure handling

  Scenario: Router returns a consistent error when the selected upstream fails
    Given router is available
    And upstream A and upstream B are configured for chat completions
    And upstream A is failing chat completions with status 503
    When send a POST request to "/v1/chat/completions" with body:
      """
      {
        "model": "auto",
        "messages": [{"role": "user", "content": "hello"}]
      }
      """
    Then response status should be 502
    And response should be valid JSON
    And response body should describe an upstream error