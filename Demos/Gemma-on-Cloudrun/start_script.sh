#!/bin/bash

set -e

# Handle shutdown signals
cleanup() {
  echo "Caught signal, shutting down..."
  kill $OLLAMA_PID $APP_PID 2>/dev/null
  wait $OLLAMA_PID $APP_PID 2>/dev/null
}

trap 'cleanup; exit 0' SIGTERM SIGINT

echo "Starting Ollama server in background..."
/bin/ollama serve &
OLLAMA_PID=$!

MAX_RETRIES=50
RETRY_INTERVAL=2
retries=0

echo "Waiting for Ollama to be ready..."
until [[ "$retries" -ge "$MAX_RETRIES" ]] || curl -s --fail http://localhost:3000; do
  sleep "$RETRY_INTERVAL"
  retries=$((retries + 1))
done

if [[ "$retries" -ge "$MAX_RETRIES" ]]; then
  echo "Ollama did not become ready within the timeout period. Exiting."
  kill $OLLAMA_PID 2>/dev/null
  exit 1
fi

echo "Ollama is ready. Starting app server..."
/app/server &
APP_PID=$!

# Wait for either process to exit
wait -n $OLLAMA_PID $APP_PID
EXIT_CODE=$?

echo "One of the processes exited with code $EXIT_CODE. Shutting down."
cleanup
exit $EXIT_CODE