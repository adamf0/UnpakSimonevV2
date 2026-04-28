# =========================
# Build Stage
# =========================
FROM registry.access.redhat.com/ubi9/go-toolset:latest AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p ./out

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -v -trimpath -ldflags="-s -w" -o ./out/app ./main.go

# =========================
# Runtime Stage
# =========================
FROM registry.access.redhat.com/ubi9/ubi-micro:latest

COPY --from=builder /usr/share/zoneinfo/Asia/Jakarta /usr/share/zoneinfo/Asia/Jakarta
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

ENV TZ=Asia/Jakarta \
    HOME=/nonexistent \
    PATH=/app

# RUN mkdir -p /app /tmp && \
#     echo 'user:x:10001:10001::/nonexistent:/sbin/nologin' > /etc/passwd && \
#     echo 'user:x:10001:' > /etc/group && \
#     chmod 1777 /tmp

COPY --from=builder --chown=10001:10001 /src/out/app /app/app

WORKDIR /app

USER 10001:10001

EXPOSE 3000

ENTRYPOINT ["/app/app"]