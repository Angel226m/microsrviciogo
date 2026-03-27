# ═══════════════════════════════════════════════════════════════
# API Gateway – Multi-stage Docker build
# ═══════════════════════════════════════════════════════════════
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=1.0.0" \
    -o /bin/api-gateway ./cmd/main.go

# ── Production image ─────────────────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /bin/api-gateway .

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

ENTRYPOINT ["./api-gateway"]
