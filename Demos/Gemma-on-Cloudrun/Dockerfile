# === Stage 1: Build proxy server ===
FROM golang:1.24-alpine as proxy_builder
LABEL stage="proxy-builder"

# Create and change to the app directory.
WORKDIR /app

# Copy local code to the container image.
COPY . ./
RUN go mod tidy
RUN go mod download

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

# Define the build argument for the model
ARG MODEL

# Store the model weights in the container image
ENV MODEL=${MODEL}


# Store model weight files in /models
ENV OLLAMA_MODELS /models
RUN /bin/ollama serve & sleep 5 && ollama pull $MODEL

WORKDIR /script

# Copy the start script to the container image.
COPY start_script.sh ./

RUN chmod +x ./start_script.sh
ENTRYPOINT ["./start_script.sh"]
