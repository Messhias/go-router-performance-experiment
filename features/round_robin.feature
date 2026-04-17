Feature: Round-robin load balancing

  Scenario: Sequential requests alternate between two upstreams
    Given router is available
    And upstream A and upstream B are configured for chat completions
    When I send 4 sequential POST requests to "/v1/chat/completions" with body:
      """
      {
        "model": "auto",
        "messages": [{"role": "user", "content": "ping"}]
      }
      """
    Then upstream handling order should be "A,B,A,B"