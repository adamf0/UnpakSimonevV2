# =========================
# Build Stage
# =========================
FROM registry.access.redhat.com/ubi9/go-toolset:latest AS builder

WORKDIR /src

# dependency cache
COPY go.mod go.sum ./
RUN go mod download

# source code
COPY . .

# prepare output folder
RUN mkdir -p /out

# build binary (gunakan main.go di root)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -v -trimpath -ldflags="-s -w" -o /out/app ./main.go

# =========================
# Runtime Stage (Hardened)
# =========================
FROM registry.access.redhat.com/ubi9/ubi-micro:latest

# timezone + ssl cert
COPY --from=builder /usr/share/zoneinfo/Asia/Jakarta /usr/share/zoneinfo/Asia/Jakarta
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

ENV TZ=Asia/Jakarta \
    HOME=/nonexistent \
    PATH=/app \
    GODEBUG=madvdontneed=1 \
    GIN_MODE=release

# create unprivileged user manually
RUN mkdir -p /app /tmp && \
    echo 'user:x:10001:10001::/nonexistent:/sbin/nologin' > /etc/passwd && \
    echo 'user:x:10001:' > /etc/group && \
    chmod 1777 /tmp

# copy binary only
COPY --from=builder /out/app /app/app

# secure permission
RUN chown -R 10001:10001 /app && \
    chmod 0555 /app/app && \
    chmod 0555 /app

# remove common abuse tools if exists
RUN rm -f \
    /bin/sh \
    /bin/bash \
    /usr/bin/bash \
    /usr/bin/sh \
    /usr/bin/curl \
    /usr/bin/wget \
    /usr/bin/nc \
    /usr/bin/netcat \
    /usr/bin/python \
    /usr/bin/python3 \
    /usr/bin/perl \
    /usr/bin/ruby \
    /usr/bin/php \
    /usr/bin/lua 2>/dev/null || true

# run as non-root
USER 10001:10001

WORKDIR /app

EXPOSE 3000/tcp

ENTRYPOINT ["/app/app"]