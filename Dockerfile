# ── Build ───────────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

WORKDIR /src

# Lapisan dependensi dipisah agar cache tidak batal setiap kali kode berubah.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
ARG COMMIT=none
ARG BUILD_DATE=unknown

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
    -ldflags "-s -w \
      -X github.com/vyerru/lapor-kampus/config.Version=${VERSION} \
      -X github.com/vyerru/lapor-kampus/config.Commit=${COMMIT} \
      -X github.com/vyerru/lapor-kampus/config.BuildDate=${BUILD_DATE}" \
    -o /out/lapor-kampus .

# ── Runtime ─────────────────────────────────────────────────────────────────
FROM alpine:3.24

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/lapor-kampus /app/lapor-kampus

# Folder log harus dapat ditulis oleh user non-root.
RUN mkdir -p /app/logs && chown -R app:app /app

USER app

ENV APP_PORT=3000 \
    LOG_LEVEL=info

EXPOSE 3000

HEALTHCHECK --interval=15s --timeout=3s --start-period=10s --retries=3 \
    CMD wget -qO- "http://127.0.0.1:${APP_PORT}/api/v1/health" >/dev/null || exit 1

ENTRYPOINT ["/app/lapor-kampus"]
