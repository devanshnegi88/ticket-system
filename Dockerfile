# ---- Build stage ----
FROM golang:1.22-bookworm AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -o /app/ticket-system ./cmd/server

# ---- Final stage ----
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    libsqlite3-0 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary
COPY --from=builder /app/ticket-system .

# Copy frontend
COPY --from=builder /app/web ./web

EXPOSE 8080

ENTRYPOINT ["./ticket-system"]