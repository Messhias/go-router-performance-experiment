Feature: Light concurrent load

  @load-light
  Scenario: Parallel chat completions stay balanced across two upstreams
    Given router is available
    And upstream A and upstream B are configured for chat completions
    When send 200 parallel POST requests to "/v1/chat/completions" with body:
      """
      {
        "model": "auto",
        "messages": [{"role": "user", "content": "ping"}]
      }
      """
    Then upstream A and upstream B should each handle between 40 and 60 percent of requests