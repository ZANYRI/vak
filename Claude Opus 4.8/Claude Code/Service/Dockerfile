# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS build

WORKDIR /src

# Cache dependencies first.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build all three binaries, statically linked.
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -ldflags="-s -w" -o /out/api ./cmd/api && \
    go build -ldflags="-s -w" -o /out/worker ./cmd/worker && \
    go build -ldflags="-s -w" -o /out/scheduler ./cmd/scheduler

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && adduser -S -G app app

WORKDIR /app
COPY --from=build /out/api /app/api
COPY --from=build /out/worker /app/worker
COPY --from=build /out/scheduler /app/scheduler

USER app
EXPOSE 8080

# Default to the API; compose overrides the command for worker/scheduler.
# Use CMD (not ENTRYPOINT) so `command:` in compose fully replaces it instead
# of being appended as arguments.
CMD ["/app/api"]
