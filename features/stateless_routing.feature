Feature: Stateless routing

  Scenario: Different clients do not get per-client affinity in load balancing
    Given router is available
    And upstream A and upstream B are configured for chat completions
    When client "alice" sends a POST request to "/v1/chat/completions" with header "X-Client-Id" "alice" and body:
      """
      {
        "model": "auto",
        "messages": [{"role": "user", "content": "one"}]
      }
      """
    And client "bob" sends a POST request to "/v1/chat/completions" with header "X-Client-Id" "bob" and body:
      """
      {
        "model": "auto",
        "messages": [{"role": "user", "content": "two"}]
      }
      """
    And client "alice" sends a POST request to "/v1/chat/completions" with header "X-Client-Id" "alice" and body:
      """
      {
        "model": "auto",
        "messages": [{"role": "user", "content": "three"}]
      }
      """
    Then upstream handling order for the last 3 requests should be "A,B,A"