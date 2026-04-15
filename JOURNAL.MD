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

After setup the the [docker-compose.yml](docker-cmopose.yml) file and test a curl by terminal I got the answers from model:

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
