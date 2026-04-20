Feature: Stateless routing

  Scenario: Different clients do not get per-client affinity in load balancing
    Given router is available
    And upstream A and upstream B are configured for chat completions
    When following POST requests are sent to "/v1/chat/completions" in order with header "X-Client-Id" and JSON bodies built from:
      | x_client_id | message_content |
      | alice       | one             |
      | bob         | two             |
      | alice       | three           |
    Then upstream handling order for the last 3 requests should be "A,B,A"