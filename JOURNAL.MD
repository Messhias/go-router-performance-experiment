## Dev journal

### 15/04/2026

I start by thinking to create a setup for docker, the reason to use it is to avoid for me install anything native in my machine for this tests.

I know that's some consideration of the loss of perfomance because I have a docker proxying it for me, but I do not want to install things in my machine unnecessary (Llama, golang) and docker keep everything working in all computers (avoid the famous phrase "in my computer it does work").

So the idea is put the localhost to work with the URL's since it will be a restful api.

I will going to use the ports 8081 and 8082 to run the models and leave the 8080 to the Gin backend

The idea is create a docker-compose.yml file with the services where I will load the `.gguf` file (I need to check where I can get one, probably in hughingfaces?).

I downloaded [this .gguf](https://huggingface.co/Qwen/Qwen2.5-0.5B-Instruct-GGUF?show_file_info=qwen2.5-0.5b-instruct-q2_k.gguf) for the *qwen* file model because it's lower since we are not aiming for quality in response, just testing the concurrency.


I downloaded [this .gguf](https://huggingface.co/ggml-org/gemma-3-1b-it-GGUF?show_file_info=gemma-3-1b-it-Q4_K_M.gguf) for the *gemma* file model because it's lower since we are not aiming for quality in response, just testing the concurrency.

I decided to `LLAMA_ARG_CTX_SIZE` to `4096` to have a not a lower larger context window and not to shorter, so with `4096` I believe we can achieve a nice value.

After setup the [docker-compose.yml](https://github.com/Messhias/go-router-performance-experiment/blob/main/docker-compose.yml) file and test a curl by terminal I got the answers from model:

Qwen:
```
curl -sS -X POST http://localhost:8082/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"gemma","messages":[{"role":"user","content":"Hello"}]}'
{"choices":[{"finish_reason":"stop","index":0,"message":{"role":"assistant","content":"Hello! How can I assist you today?"}}],"created":1776268502,"model":"qwen2.5-0.5b-instruct-q2_k.gguf","system_fingerprint":"b8783-e21cdc11a","object":"chat.completion","usage":{"completion_tokens":10,"prompt_tokens":30,"total_tokens":40,"prompt_tokens_details":{"cached_tokens":0}},"id":"chatcmpl-ITNgzCzpMJBP9ObgozDvsi2olmBEE1CK","timings":{"cache_n":0,"prompt_n":30,"prompt_ms":767.391,"prompt_per_token_ms":25.5797,"prompt_per_second":39.09349992376767,"predicted_n":10,"predicted_ms":390.006,"predicted_per_token_ms":39.0006,"predicted_per_second":25.640631169776874}}
```

Gemma:
```    
curl -sS -X POST http://localhost:8081/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"gemma","messages":[{"role":"user","content":"Hello"}]}'
{"choices":[{"finish_reason":"stop","index":0,"message":{"role":"assistant","content":"Hello there! How can I help you today? 😊 \n\nDo you have any questions for me, or would you like to chat about something?"}}],"created":1776268470,"model":"gemma-3-1b-it-Q4_K_M.gguf","system_fingerprint":"b8783-e21cdc11a","object":"chat.completion","usage":{"completion_tokens":31,"prompt_tokens":10,"total_tokens":41,"prompt_tokens_details":{"cached_tokens":0}},"id":"chatcmpl-s83tKXpmOfdMtCSEB0qr6zbdFqzZh0Kq","timings":{"cache_n":0,"prompt_n":10,"prompt_ms":602.724,"prompt_per_token_ms":60.272400000000005,"prompt_per_second":16.591341974104232,"predicted_n":31,"predicted_ms":2827.211,"predicted_per_token_ms":91.20035483870967,"predicted_per_second":10.964869618857596}}  
```

So I belive for today after my research and used the right tools I achieve to setup the Llama models.

[Here's the MR (#1) I created for this task](https://github.com/Messhias/go-router-performance-experiment/pull/1)

-----

Moving on I decided to start the setup golang essentials (and it's packages)


One decision I took in relation of it, it is not use fixed versions of the images, I always hate legacy and the promise of "later we update it" is a lie that every developer likes to tell to himself.

My philosophys is always to "legacy here no" if something brokes, we can fix it up (also fix the things in specific version means that you keeping the security issues as well).

[Here's the MR (#2) I created for this task](https://github.com/Messhias/go-router-performance-experiment/pull/2)


----

### 16/04/2026

Initial plan for BDD + TDD

*Scenario 1* -> OpenAI compatible.
*Scenario 2* -> Round robin algorithm.
*Scenario 3* -> Proxy transparency.
*Scenario 4* -> Stateless.
*Scenario 5* -> Handle upstream failures.


#### Pyramid tests

TDD UNIT:
- Upstream selection.
- road robin algorithm.
- request validation.

TDD Component:
- handler Gin 
- httptest with mocked upstream (we won't use real external network for the sake of simplicity).

#### BDD

BDD acceptance with Godog (write the step definitions + feature files).

----

For now  I am thinkin in this folder structure to work with it:

- internal/router (setup of Gin and routes)
- internal/proxy (forward request/response)
- internal/balancer (round-robin thread-safe)
- internal/upstream (backend config for config/cliente)
- features/ + *_test.go (BDD e TDD)

#### My mental model for the contracts


- Endpoint: POST /v1/chat/completions
- Content-Type JSON
- model body with messages
- response with OpenAI-compatible
- Load balance between the 2 instances.
- The api needs to be stateless


## My tasks backlog

- TKT-001 — Test Bootstrap and Conventions [DONE]


Objective: Standardize test structure (unit/component/BDD) and execution commands.
BDD (Red): Technical scenario “feature suite executes and reports failure correctly”.
TDD: Create minimal sanity test (go test runs project).
Ready when: There is a single command to run tests per layer and short documentation in the internal README.

- TKT-002 — Gherkin Scenarios from the README (source of truth) [DONE]

Objective: Transform requirements into .feature files.
BDD (Red): Create scenarios for OpenAI compatibility, round-robin, transparency, stateless, and upstream error.
TDD: n/a (executable specification ticket).

Ready when: All scenarios exist and fail due to lack of implementation (expected).

- TKT-003 — Test doubles of OpenAI-compatible upstream

Objective: Create two fake upstreams with httptest.Server.
Partial BDD (Green): Steps can send requests and capture which upstream responded.
TDD: Unit tests of the fakes (status, JSON body, upstream identification).

Ready when: Upstreams A/B simulate /v1/chat/completions deterministically.

- TKT-004 — Request parser/validator contract

Objective: Validate the minimum structure of the OpenAI chat completions request.

BDD (Red): Invalid payload scenario returns expected error.

TDD: Unit tests for valid/invalid payload.

Ready when: Validation covers main cases without coupling to the Gin handler.

- TKT-005 — Round-robin thread-safe (core)

Objective: A/B alternating selection component safe for concurrency.
BDD (Red): Sequential request scenario alternates upstream.
TDD: Unit tests for alternation and basic concurrency.

Ready when: Algorithm passes deterministic tests and does not depend on per-client state.

- TKT-006 — Gin handler for chat completions (minimum)

Objective: Expose POST /v1/chat/completions.

BDD (Red→Green): Basic success scenario.

TDD: Component test with httptest.NewRecorder for status/body.

Ready when: Endpoint responds with OpenAI-compatible JSON via fake upstream. 

- TKT-007 — Transparent Proxy for Essential Body/Headers

Objective: Preserve request/response with minimal transformation.
BDD (Red): Transparency scenario validates important fields.
TDD: Component tests for method forwarding, body, and critical headers (Content-Type, etc.).
Ready when: Payload differences are only those necessary for routing.

- TKT-008 — Upstream Fault Handling

Objective: Define and implement error behavior.

BDD (Red→Green): “upstream unavailable” scenario.

TDD: Tests for timeout, 5xx, connection refused.

Ready when: Error is consistent, observable, and does not crash the router.

- TKT-009 — Guarantee of Stateless Behavior

Objective: Prove absence of affinity/session in routing.

BDD (Red→Green): Multiple clients globally alternate through the same RR.

TDD: Component testing with different headers/client IDs without impacting selection.
Ready when: No per-user state influences load balancing.

- TKT-010 — Godog + Go Test Integration

Objective: Run BDD alongside the standard Go suite.
BDD: Execution of features within TestFeatures.

TDD: "Wiring" test of the suite.

Ready when: Go test executes unit/component/BDD in the same flow.

- TKT-011 — Concurrency and Stability Testing


Objective: Validate behavior under parallel light load.
BDD: tagged (@load-light) scenario for approximate distribution.

TDD: Concurrent testing on the balancer/handler.

Ready when: No race conditions and expected distribution between A/B.

- TKT-012 — Real smoke with llama.cpp + acceptance checklist

Objective: Validate end-to-end with the two real instances.
BDD: Mark scenarios as “accepted” with evidence (summarized commands and outputs).

TDD: n/a (real integration validation).

Ready when: curl on the router works, round-robin observable, and compatible format.



### 17/04/2026

Start the work in the TKT-001/002 tasks

MR: [TKT-001 [DONE]](https://github.com/Messhias/go-router-performance-experiment/pull/4)

MR: [TKT-002 [DONE]](https://github.com/Messhias/go-router-performance-experiment/pull/5)


### 18/04/2026

MR: [TKT-003 [DONE]](https://github.com/Messhias/go-router-performance-experiment/pull/6)

MR: [TKT-004 [DONE]](https://github.com/Messhias/go-router-performance-experiment/pull/7)

MR: [TKT-005 [DONE]](https://github.com/Messhias/go-router-performance-experiment/pull/8)

### 19/04/2026

MR: [TKT-006 [DONE]](https://github.com/Messhias/go-router-performance-experiment/pull/9)

MR: [TKT-007 [DONE]](https://github.com/Messhias/go-router-performance-experiment/pull/10)

### 20/04/2026

MR: [TKT-008 [DONE]](https://github.com/Messhias/go-router-performance-experiment/pull/11)

MR: [TKT-009 [DONE]](https://github.com/Messhias/go-router-performance-experiment/pull/12)

MR: [TKT-010 [DONE]](https://github.com/Messhias/go-router-performance-experiment/pull/12)

Evidence: 
```
go test ./... -v -race
?   	messhias/router-expirement	[no test files]
=== RUN   TestFeatures
Feature: OpenAI-compatible chat completions
=== RUN   TestFeatures/Successful_chat_completion_with_minimal_valid_payload
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /v1/chat/completions      --> messhias/router-expirement/internal/router.registerChatRoutes.func1 (1 handlers)

  Background:
    Given router is available                                              # common.go:52 -> messhias/router-expirement/tests/acceptance_tests.givenRouterIsAvailable

  Scenario: Successful chat completion with minimal valid payload          # openai_compatibility.feature:6
    Given upstream responds with an OpenAI-compatible chat completion      # common.go:54 -> messhias/router-expirement/tests/acceptance_tests.InitCommon.givenUpstreamAAndUpstreamB.func3
    When send a POST request to "/v1/chat/completions" with body:          # openai_compatibility.go:120 -> messhias/router-expirement/tests/acceptance_tests.whenPostRequest
      """
      {
        "model": "auto",
        "messages": [{"role": "user", "content": "Hello"}]
      }
      """
    Then response status should be 200                                     # openai_compatibility.go:121 -> messhias/router-expirement/tests/acceptance_tests.thenResponseStatus200
    And response should be valid JSON                                      # openai_compatibility.go:123 -> messhias/router-expirement/tests/acceptance_tests.thenShouldValidJson
    And response should contain an OpenAI-compatible chat completion shape # openai_compatibility.go:124 -> messhias/router-expirement/tests/acceptance_tests.thenShouldContainOpenAICompatible
=== RUN   TestFeatures/Invalid_chat_completion_payload_returns_a_client_error
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /v1/chat/completions      --> messhias/router-expirement/internal/router.registerChatRoutes.func1 (1 handlers)

  Scenario: Invalid chat completion payload returns a client error # openai_compatibility.feature:19
    When send a POST request to "/v1/chat/completions" with body:  # openai_compatibility.go:120 -> messhias/router-expirement/tests/acceptance_tests.whenPostRequest
      """
      {
        "model": "",
        "messages": []
      }
      """
    Then response status should be 400                             # openai_compatibility.go:122 -> messhias/router-expirement/tests/acceptance_tests.thenResponseStatus400

Feature: Transparent proxy behavior
=== RUN   TestFeatures/Essential_request_body_and_headers_reach_the_upstream_with_minimal_change
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /v1/chat/completions      --> messhias/router-expirement/internal/router.registerChatRoutes.func1 (1 handlers)

  Scenario: Essential request body and headers reach the upstream with minimal change       # proxy_transparency.feature:3
    Given router is available                                                               # common.go:52 -> messhias/router-expirement/tests/acceptance_tests.givenRouterIsAvailable
    And upstream A is configured to echo the received request for chat completions          # proxy_transparency.go:233 -> messhias/router-expirement/tests/acceptance_tests.InitProxyTransparency.func1
    When I send a POST request to "/v1/chat/completions" with headers:                      # proxy_transparency.go:236 -> messhias/router-expirement/tests/acceptance_tests.InitProxyTransparency.func2
      | name         | value            |
      | Content-Type | application/json |
    And body:                                                                               # proxy_transparency.go:239 -> messhias/router-expirement/tests/acceptance_tests.InitProxyTransparency.func3
      """
      {
        "model": "auto",
        "messages": [{"role": "user", "content": "Hello"}]
      }
      """
    Then upstream A should have received the same JSON body                                 # proxy_transparency.go:242 -> messhias/router-expirement/tests/acceptance_tests.InitProxyTransparency.func4
    And upstream A should have received header "Content-Type" with value "application/json" # proxy_transparency.go:245 -> messhias/router-expirement/tests/acceptance_tests.InitProxyTransparency.func5

Feature: Round-robin load balancing
=== RUN   TestFeatures/Sequential_requests_alternate_between_two_upstreams
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /v1/chat/completions      --> messhias/router-expirement/internal/router.registerChatRoutes.func1 (1 handlers)

  Scenario: Sequential requests alternate between two upstreams                 # round_robin.feature:3
    Given router is available                                                   # common.go:52 -> messhias/router-expirement/tests/acceptance_tests.givenRouterIsAvailable
    And upstream A and upstream B are configured for chat completions           # common.go:53 -> messhias/router-expirement/tests/acceptance_tests.InitCommon.givenUpstreamAAndUpstreamB.func2
    When I send 4 sequential POST requests to "/v1/chat/completions" with body: # round_robin.go:108 -> messhias/router-expirement/tests/acceptance_tests.whenISend4SequentialPost
      """
      {
        "model": "auto",
        "messages": [{"role": "user", "content": "ping"}]
      }
      """
    Then upstream handling order should be "A,B,A,B"                            # round_robin.go:109 -> messhias/router-expirement/tests/acceptance_tests.thenUpstreamHandlingOrderShouldBe

Feature: Stateless routing
=== RUN   TestFeatures/Different_clients_do_not_get_per-client_affinity_in_load_balancing
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /v1/chat/completions      --> messhias/router-expirement/internal/router.registerChatRoutes.func1 (1 handlers)

  Scenario: Different clients do not get per-client affinity in load balancing                                                     # stateless_routing.feature:3
    Given router is available                                                                                                      # common.go:52 -> messhias/router-expirement/tests/acceptance_tests.givenRouterIsAvailable
    And upstream A and upstream B are configured for chat completions                                                              # common.go:53 -> messhias/router-expirement/tests/acceptance_tests.InitCommon.givenUpstreamAAndUpstreamB.func2
    When following POST requests are sent to "/v1/chat/completions" in order with header "X-Client-Id" and JSON bodies built from: # stateless_routing.go:97 -> messhias/router-expirement/tests/acceptance_tests.whenClientSendsPOSTWithHeaderAndBody
      | x_client_id | message_content |
      | alice       | one             |
      | bob         | two             |
      | alice       | three           |
    Then upstream handling order for the last 3 requests should be "A,B,A"                                                         # stateless_routing.go:98 -> messhias/router-expirement/tests/acceptance_tests.thenUpstreamHandlingOrderLast3ABABA

Feature: Upstream failure handling
=== RUN   TestFeatures/Router_returns_a_consistent_error_when_the_selected_upstream_fails
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /v1/chat/completions      --> messhias/router-expirement/internal/router.registerChatRoutes.func1 (1 handlers)

  Scenario: Router returns a consistent error when the selected upstream fails # upstream_failures.feature:3
    Given router is available                                                  # common.go:52 -> messhias/router-expirement/tests/acceptance_tests.givenRouterIsAvailable
    And upstream A and upstream B are configured for chat completions          # common.go:53 -> messhias/router-expirement/tests/acceptance_tests.InitCommon.givenUpstreamAAndUpstreamB.func2
    And upstream A is failing chat completions with status 503                 # upstream_failures.go:114 -> messhias/router-expirement/tests/acceptance_tests.givenUpstreamAFailingChatCompletions503
    When send a POST request to "/v1/chat/completions" with body:              # openai_compatibility.go:120 -> messhias/router-expirement/tests/acceptance_tests.whenPostRequest
      """
      {
        "model": "auto",
        "messages": [{"role": "user", "content": "hello"}]
      }
      """
    Then response status should be 502                                         # upstream_failures.go:115 -> messhias/router-expirement/tests/acceptance_tests.thenTheResponseStatusShouldBe
    And response should be valid JSON                                          # openai_compatibility.go:123 -> messhias/router-expirement/tests/acceptance_tests.thenShouldValidJson
    And response body should describe an upstream error                        # upstream_failures.go:117 -> messhias/router-expirement/tests/acceptance_tests.thenResponseBodyDescribesUpstreamError

6 scenarios (6 passed)
30 steps (30 passed)
36.906458ms
--- PASS: TestFeatures (0.04s)
    --- PASS: TestFeatures/Successful_chat_completion_with_minimal_valid_payload (0.01s)
    --- PASS: TestFeatures/Invalid_chat_completion_payload_returns_a_client_error (0.00s)
    --- PASS: TestFeatures/Essential_request_body_and_headers_reach_the_upstream_with_minimal_change (0.00s)
    --- PASS: TestFeatures/Sequential_requests_alternate_between_two_upstreams (0.00s)
    --- PASS: TestFeatures/Different_clients_do_not_get_per-client_affinity_in_load_balancing (0.00s)
    --- PASS: TestFeatures/Router_returns_a_consistent_error_when_the_selected_upstream_fails (0.00s)
PASS
ok  	messhias/router-expirement/features	1.857s
=== RUN   TestChatRequestDto_ShouldPass
=== RUN   TestChatRequestDto_ShouldPass/valid_minimal
=== RUN   TestChatRequestDto_ShouldPass/invalid_json
=== RUN   TestChatRequestDto_ShouldPass/empty_model_and_empty_message_fields
=== RUN   TestChatRequestDto_ShouldPass/empty_model_empty_content
=== RUN   TestChatRequestDto_ShouldPass/empty_model_empty_role
--- PASS: TestChatRequestDto_ShouldPass (0.00s)
    --- PASS: TestChatRequestDto_ShouldPass/valid_minimal (0.00s)
    --- PASS: TestChatRequestDto_ShouldPass/invalid_json (0.00s)
    --- PASS: TestChatRequestDto_ShouldPass/empty_model_and_empty_message_fields (0.00s)
    --- PASS: TestChatRequestDto_ShouldPass/empty_model_empty_content (0.00s)
    --- PASS: TestChatRequestDto_ShouldPass/empty_model_empty_role (0.00s)
PASS
ok  	messhias/router-expirement/internal/DTO	(cached)
?   	messhias/router-expirement/internal/acceptance	[no test files]
=== RUN   TestInvalidBalancerCreation_ShouldFail
--- PASS: TestInvalidBalancerCreation_ShouldFail (0.00s)
=== RUN   TestBalancerAlternates_ShouldPass
--- PASS: TestBalancerAlternates_ShouldPass (0.00s)
PASS
ok  	messhias/router-expirement/internal/balancer	(cached)
?   	messhias/router-expirement/internal/config	[no test files]
?   	messhias/router-expirement/internal/proxy	[no test files]
=== RUN   TestPOSTChatCompletions_proxiesToUpstream_openAICompatibleJSON_ShouldPass
--- PASS: TestPOSTChatCompletions_proxiesToUpstream_openAICompatibleJSON_ShouldPass (0.00s)
=== RUN   TestPOSTChatCompletions_upstreamReceivesPostMethod_ShouldPass
--- PASS: TestPOSTChatCompletions_upstreamReceivesPostMethod_ShouldPass (0.00s)
=== RUN   TestPOSTChatCompletions_upstreamReceivesSameBodyAndContentType_ShouldPass
--- PASS: TestPOSTChatCompletions_upstreamReceivesSameBodyAndContentType_ShouldPass (0.00s)
=== RUN   TestTimeout_ShouldPass
--- PASS: TestTimeout_ShouldPass (0.10s)
PASS
ok  	messhias/router-expirement/internal/router	(cached)
=== RUN   TestUpstreams_ShouldPass
--- PASS: TestUpstreams_ShouldPass (0.01s)
=== RUN   TestWrongMethodInWrongUrl_ShouldFail
--- PASS: TestWrongMethodInWrongUrl_ShouldFail (0.00s)
PASS
ok  	messhias/router-expirement/internal/upstreamfake	(cached)
?   	messhias/router-expirement/tests/acceptance_tests	[no test files]
```