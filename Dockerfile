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

# buat identity user manual + shell disable total
# /sbin/nologin dipakai jika ada attempt login
COPY <<EOF /etc/passwd
root:x:0:0:root:/root:/sbin/nologin
user:x:10001:10001:App User:/nonexistent:/sbin/nologin
EOF

COPY <<EOF /etc/group
root:x:0:
user:x:10001:
EOF

# block shell path tambahan (defense in depth)
COPY <<EOF /sbin/nologin
#!/bin/sh
echo "This account is not available."
exit 1
EOF

COPY --from=builder --chown=10001:10001 /src/out/app /app/app

WORKDIR /app

USER 10001:10001

EXPOSE 3000

ENTRYPOINT ["/app/app"]