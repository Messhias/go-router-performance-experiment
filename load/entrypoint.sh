#!/bin/sh
set -eu

ARTILLERY_VERSION="${ARTILLERY_VERSION:-latest}"
export BUN_INSTALL_CACHE_DIR="${BUN_INSTALL_CACHE_DIR:-/work/out/.buncache}"

apk add --no-cache curl >/dev/null

ready=0
i=0
while [ "$i" -lt 60 ]; do
  i=$((i + 1))
  if curl -sf --connect-timeout 5 --max-time 8 "http://llama-gemma:8080/v1/models" >/dev/null 2>&1 \
    && curl -sf --connect-timeout 5 --max-time 8 "http://llama-qwen:8080/v1/models" >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 5
done

if [ "$ready" -ne 1 ]; then
  echo "llama.cpp readiness timeout" >&2
  exit 1
fi

mkdir -p /work/out/.buncache

RUN_JSON="/work/out/run.json"
HTML_OUT="/work/out/report.html"

bunx --bun "artillery@${ARTILLERY_VERSION}" run --output "$RUN_JSON" /work/artillery.yml
bunx --bun artillery@2.0.21 report "$RUN_JSON" -o "$HTML_OUT"
echo "$HTML_OUT"
