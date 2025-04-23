# === Stage 1: Build proxy server ===
FROM golang:1.24-alpine as proxy_builder
LABEL stage="proxy-builder"

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -v -o server


# === Stage 2: Prepare Ollama and Download Model ===
# Use the official ollama image as a base to get the binary and pull the model

FROM ollama/ollama:latest as ollama_builder
LABEL stage="ollama-builder"

# Install curl. It's used in the start_script.sh to check Ollama health status.
RUN apt-get update && apt-get install -y curl

# Copy the proxy server binary to the ollama image.
COPY --from=proxy_builder /app/server /app/server

# Listen to localhost, port 3000
ENV OLLAMA_HOST http://localhost:3000

# Reduce logging verbosity
ENV OLLAMA_DEBUG false

# Never unload model weights from the GPU
ENV OLLAMA_KEEP_ALIVE -1

ENV OLLAMA_ORIGINS *
# Store the model weights in the container image
ENV MODEL=gemma3:4b

# Store model weight files in /models
ENV OLLAMA_MODELS /models
RUN /bin/ollama serve & sleep 5 && ollama pull $MODEL

WORKDIR /script

# Copy the start script to the container image.
COPY start_script.sh ./

RUN chmod +x ./start_script.sh
ENTRYPOINT ["./start_script.sh"]
