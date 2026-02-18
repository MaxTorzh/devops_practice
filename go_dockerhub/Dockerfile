FROM golang:1.22-alpine AS builder

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="\
    -X 'github.com/MaxTorzh/go_dockerhub/internal/version.Version=${VERSION}' \
    -X 'github.com/MaxTorzh/go_dockerhub/internal/version.Commit=${COMMIT}' \
    -X 'github.com/MaxTorzh/go_dockerhub/internal/version.BuildTime=${BUILD_TIME}'" \
    -o /app/bin/api ./cmd/api

FROM alpine:3.19

LABEL org.opencontainers.image.source="https://github.com/MaxTorzh/go-dockerhub"
LABEL org.opencontainers.image.description="Demo app for Docker Hub"
LABEL org.opencontainers.image.licenses="MIT"

RUN apk --no-cache add ca-certificates tzdata

RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup /app/bin/api /app/api

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -q --spider http://localhost:8080/health || exit 1

CMD ["/app/api"]