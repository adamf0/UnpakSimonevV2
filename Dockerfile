# =========================
# Build Stage
# =========================
FROM golang:1.25.4-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /out

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/app \
    ./main.go

# =========================
# Runtime Stage
# =========================
FROM registry.access.redhat.com/ubi9/ubi-micro:latest

# timezone
#COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# ssl cert
#COPY --from=builder /etc/ssl/certs /etc/ssl/certs

ENV TZ=Asia/Jakarta

# user manual
COPY <<EOF /etc/passwd
appuser:x:10001:10001:Application User:/nonexistent:/sbin/nologin
EOF

COPY <<EOF /etc/group
appuser:x:10001:
EOF

# binary
COPY --from=builder --chown=10001:10001 /out/app /app/app

WORKDIR /app

USER 10001:10001

EXPOSE 3000

ENTRYPOINT ["/app/app"]
