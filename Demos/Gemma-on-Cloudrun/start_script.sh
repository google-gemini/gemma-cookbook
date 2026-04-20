#!/bin/bash

set -e
echo "Starting Ollama server in background..."
/bin/ollama serve &
OLLAMA_PID=$! # Get PID of the background process


MAX_RETRIES=50 # Wait for a maximum of 30 attempts
RETRY_INTERVAL=2 # Wait 2 seconds between attempts
retries=0

echo "Waiting for Ollama to be ready..."

until [[ "$retries" -ge "$MAX_RETRIES" ]] || curl -s --fail http://localhost:3000; do
  echo "Attempt $retries: Waiting for Ollama to be ready..."
  sleep "$RETRY_INTERVAL"
  retries=$((retries + 1))
done

if [[ "$retries" -ge "$MAX_RETRIES" ]]; then
  echo "Ollama did not become ready within the timeout period. Exiting."
  exit 1 # Or some other error code
else
  echo "Ollama is ready."
  # Proceed with the rest of your script
fi

echo "Starting app server in foreground..."
/app/server &
APP_PID=$! # Get PID of the background process

# Wait for either process to exit
wait -n $OLLAMA_PID $APP_PID

# Optional: Add trap for signal handling (SIGTERM from docker stop)
trap "echo 'Caught SIGTERM, shutting down...'; kill $OLLAMA_PID $APP_PID; wait $OLLAMA_PID; wait $APP_PID" SIGTERM SIGINT

echo "One of the processes exited."
# Wait for remaining process if necessary (or kill them via trap)
wait $OLLAMA_PID
wait $APP_PID