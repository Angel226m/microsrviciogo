FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/product-service ./cmd/main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata curl
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /bin/product-service .
USER appuser
EXPOSE 8082
HEALTHCHECK --interval=10s --timeout=3s --retries=3 CMD curl -f http://localhost:8082/health || exit 1
ENTRYPOINT ["./product-service"]
