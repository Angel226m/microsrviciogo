FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/notification-service ./cmd

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 1001 appuser
COPY --from=builder /bin/notification-service /usr/local/bin/
USER appuser
EXPOSE 8086
ENTRYPOINT ["notification-service"]
