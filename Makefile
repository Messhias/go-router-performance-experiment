.PHONY: full-smoke-test start smoke-test artillery

start:
	@cp .env.example .env
	@ids=$$(docker ps -aq); \
	if [ -n "$$ids" ]; then \
		docker stop $$ids; \
	else \
		echo "No running containers found."; \
	fi
	@docker compose up -d --build
	@go test ./...
	@echo "Waiting for llama.cpp upstreams (GET /v1/models on 8081 and 8082)..."
	@echo "(Up to ~5 min: 60 attempts x 5s; each check uses connect+transfer max 8s.)"
	@ready=0; \
	for i in $$(seq 1 60); do \
	  if curl -sf --connect-timeout 5 --max-time 8 http://localhost:8081/v1/models >/dev/null 2>&1 \
	     && curl -sf --connect-timeout 5 --max-time 8 http://localhost:8082/v1/models >/dev/null 2>&1; then \
	    ready=1; echo "Upstreams ready."; break; \
	  fi; \
	  echo "  attempt $$i/60 ..."; sleep 5; \
	done; \
	if [ "$$ready" != 1 ]; then \
	  echo "ERROR: llama.cpp services did not become ready in time." >&2; \
	  exit 1; \
	fi

full-smoke-test:
	$(MAKE) start
	$(MAKE) smoke-test


smoke-test:
	@echo ""
	@echo "Smoke: POST http://localhost:8080/v1/chat/completions"
	@echo "If nothing prints for a while: CPU inference is still running (max 360s)."
	@echo "If it fails fast: run 'make start' first or check docker compose logs."
	@echo ""
	@curl -fS --connect-timeout 15 --max-time 360 -X POST http://localhost:8080/v1/chat/completions \
		-H "Content-Type: application/json" \
		-d '{"model":"auto","messages":[{"role":"user","content":"Hi, what you can do for me?"}],"max_tokens":32}'
	@echo ""


artillery:
	docker compose --profile load run --rm artillery
