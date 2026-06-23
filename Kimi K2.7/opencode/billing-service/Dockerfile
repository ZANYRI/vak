# Build stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /build/bin/api ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /build/bin/worker ./cmd/worker && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /build/bin/scheduler ./cmd/scheduler

# Runtime stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /build/bin/api /app/api
COPY --from=builder /build/bin/worker /app/worker
COPY --from=builder /build/bin/scheduler /app/scheduler
COPY --from=builder /build/migrations /app/migrations
COPY --from=builder /build/docs /app/docs

RUN addgroup -S billing && adduser -S billing -G billing
USER billing

EXPOSE 8080 9090
