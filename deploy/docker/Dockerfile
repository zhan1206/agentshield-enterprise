# Build stage
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/ .
RUN go mod tidy && go build -o agentshield-server ./cmd/server

FROM rust:1.74-alpine AS sandbox-builder
WORKDIR /app/sandbox
COPY sandbox/ .
RUN cargo build --release

# Runtime stage
FROM python:3.11-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=backend-builder /app/backend/agentshield-server /usr/local/bin/
COPY --from=sandbox-builder /app/sandbox/target/release/libagentshield_sandbox.so /usr/local/lib/

WORKDIR /app
COPY ai/requirements.txt ai/requirements.txt
RUN pip install --no-cache-dir -r ai/requirements.txt
COPY ai/ ai/
COPY configs/ configs/

ENV CONFIG_PATH=/app/configs/config.yaml
EXPOSE 8080 8001

CMD ["agentshield-server"]
