# === Stage 1: Build proxy server ===
FROM golang:1.24-alpine AS proxy_builder

# Create and change to the app directory.
WORKDIR /app

# Copy dependency files first to leverage Docker layer caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the local code and build.
COPY . ./
RUN CGO_ENABLED=0 go build -v -o server


# === Stage 2: Prepare Ollama and Download Model ===
FROM ollama/ollama:0.16.1 AS model_downloader

# Define the build argument for the model
ARG MODEL
ENV MODEL=${MODEL}
ENV OLLAMA_MODELS=/models

# Store model weight files in /models
RUN /bin/ollama serve & sleep 5 && ollama pull "$MODEL"


# === Stage 3: Final Runtime Image ===
FROM ollama/ollama:0.16.1

# Install curl for health checks and ca-certificates for secure connections.
# Use --no-install-recommends to keep the image small and clean up apt cache to reduce vulnerabilities.
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user and group for better security.
RUN groupadd -r ollama && useradd -r -g ollama -d /home/ollama -m ollama

# Create necessary directories and set permissions.
RUN mkdir -p /models /app && \
    chown -R ollama:ollama /models /app

# Copy binaries and model files from previous stages.
COPY --from=proxy_builder /app/server /app/server
COPY --from=model_downloader /models /models
COPY start_script.sh /app/start_script.sh

# Ensure the start script is executable and owned by the non-root user.
RUN chmod +x /app/start_script.sh && \
    chown ollama:ollama /app/server /app/start_script.sh

# Environment variables for Ollama.
ENV OLLAMA_HOST=http://localhost:3000
ENV OLLAMA_MODELS=/models
ENV OLLAMA_DEBUG=false
ENV OLLAMA_KEEP_ALIVE=-1
ENV OLLAMA_ORIGINS=*
ENV HOME=/home/ollama

# Switch to the non-root user.
USER ollama
WORKDIR /app

ENTRYPOINT ["./start_script.sh"]
