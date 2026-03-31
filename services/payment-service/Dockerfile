# ═══════════════════════════════════════════════════════════════
# Servicio de Pagos – Build multi-etapa con seguridad OWASP
# ═══════════════════════════════════════════════════════════════
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=1.0.0 -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o /bin/payment-service ./cmd/main.go

# ── Imagen de producción mínima ──────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl \
    && apk upgrade --no-cache \
    && addgroup -S appgroup && adduser -S appuser -G appgroup \
    && mkdir -p /app && chown appuser:appgroup /app

WORKDIR /app

COPY --from=builder /bin/payment-service .

USER appuser

EXPOSE 8084

HEALTHCHECK --interval=10s --timeout=3s --retries=3 \
    CMD curl -f http://localhost:8084/health || exit 1

ENTRYPOINT ["./payment-service"]
