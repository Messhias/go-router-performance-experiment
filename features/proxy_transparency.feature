Feature: Transparent proxy behavior

  Scenario: Essential request body and headers reach the upstream with minimal change
    Given router is available
    And upstream A is configured to echo the received request for chat completions
    When I send a POST request to "/v1/chat/completions" with headers:
      | name           | value            |
      | Content-Type   | application/json |
    And body:
      """
      {
        "model": "auto",
        "messages": [{"role": "user", "content": "Hello"}]
      }
      """
    Then upstream A should have received the same JSON body
    And upstream A should have received header "Content-Type" with value "application/json"