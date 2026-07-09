# ---- Build stage ----
FROM golang:1.22-bookworm AS builder

WORKDIR /app

# Cache dependency downloads separately from source changes.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO is required by the mattn/go-sqlite3 driver used under gorm.io/driver/sqlite.
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/ticket-system ./cmd/server

# ---- Final stage ----
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    libsqlite3-0 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/ticket-system .

EXPOSE 8080

ENTRYPOINT ["./ticket-system"]
