# Go Router Performance Experiment
Overview
Implement a stateless REST API microservice in Go to evaluate potential performance improvements in concurrency handling and HTTP request routing compared to our existing Python-based stack. The objective is to explore where a lightweight Go-based routing layer could improve performance characteristics while remaining compatible with OpenAI-style APIs. 
Task

1. Deploy Two llama.cpp OpenAI-Compatible API Instances
- Run two independent instances of the llama.cpp OpenAI-compatible server.
- Models:
- - ggml-org/gemma-3-1b-it-GGUF
- - qwen2.5-0.5b (GGUF format)
- Requirements:
- - Each model must run as a separate HTTP endpoint.
- - Both instances must expose OpenAI-compatible REST APIs.
- - CPU-only execution (no CUDA required).
- - Example request (must work via curl):
```
curl -X POST http://localhost:<PORT>/v1/chat/completions \
-H "Content-Type: application/json" \
-d '{
  "model": "gemma",
  "messages": [
    {"role": "user", "content": "Hello"}
  ]
}'
```

2. Implement Go Router (Gin)
Create a Go-based router using Gin that:

- Accepts OpenAI-compatible chat completion requests.

- Routes incoming requests across both llama.cpp instances using round-robin load balancing.

- Maintains stateless operation.


Requirements:

- Inputs must be curl-based requests.


- Response format must remain OpenAI-compatible.


- Router should behave as a transparent proxy with minimal transformation.


Example router request:
```
curl -X POST http://localhost:8080/v1/chat/completions \
-H "Content-Type: application/json" \
-d '{
  "model": "auto",
  "messages": [
    {"role": "user", "content": "Test routing"}
  ]
}'
```
