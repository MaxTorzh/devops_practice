FROM golang:1.22-alpine AS builder

ARG SERVICE

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/service ./cmd/${SERVICE}

FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

WORKDIR /app

COPY --from=builder /app/service /app/service

EXPOSE 8080

HEALTHCHECK --interval=5s --timeout=3s --start-period=2s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ping || exit 1

CMD ["/app/service"]