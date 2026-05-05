# Build stage - Go backend
FROM golang:1.21-alpine AS backend-builder
RUN apk add --no-cache git gcc musl-dev
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 go build -o agentshield-server ./cmd/server

# Build stage - Rust sandbox engine
FROM rust:1.74-bookworm AS sandbox-builder
WORKDIR /app/sandbox
COPY sandbox/Cargo.toml sandbox/Cargo.lock ./
RUN mkdir src && echo "" > src/lib.rs && cargo build --release 2>/dev/null || true
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
