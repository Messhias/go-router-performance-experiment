Feature: OpenAI-compatible chat completions

  Scenario: Successful chat completion with minimal valid payload
    Given router is available
    And upstream responds with an OpenAI-compatible chat completion
    When send a POST request to "/v1/chat/completions" with body:
      """
      {
        "model": "auto",
        "messages": [{"role": "user", "content": "Hello"}]
      }
      """
    Then response status should be 200
    And response should be valid JSON
    And response should contain an OpenAI-compatible chat completion shape