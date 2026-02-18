FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates

ENV APP_NAME="DockerEnvApp" \
    APP_ENV="production" \
    PORT=8080 \
    LOG_LEVEL="warn" \
    MAX_WORKERS=20 \
    TIMEOUT="30s" \
    FEATURE_A="false" \
    FEATURE_B="false" \
    FEATURE_C="false"

ENV DATABASE_URL="postgres://localhost:5432/app" \
    REDIS_URL="redis://localhost:6379"

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -q --spider http://localhost:8080/health || exit 1

CMD ["./main"]