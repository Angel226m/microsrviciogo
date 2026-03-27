FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/order-service ./cmd/main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates curl
RUN addgroup -S app && adduser -S app -G app
WORKDIR /app
COPY --from=builder /bin/order-service .
USER app
EXPOSE 8083
HEALTHCHECK --interval=10s --timeout=3s --retries=3 CMD curl -f http://localhost:8083/health || exit 1
ENTRYPOINT ["./order-service"]
