# =========================
# Build Stage
# =========================
FROM registry.access.redhat.com/ubi9/go-toolset:latest AS builder

WORKDIR /src

# copy dependency file only first (cache layer)
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/app .

# =========================
# Runtime Stage (Hardened)
# =========================
FROM registry.access.redhat.com/ubi9/ubi-micro:latest

# timezone + cert only
COPY --from=builder /usr/share/zoneinfo/Asia/Jakarta /usr/share/zoneinfo/Asia/Jakarta
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

ENV TZ=Asia/Jakarta \
    HOME=/nonexistent \
    PATH=/app \
    GODEBUG=madvdontneed=1

# create minimal passwd/group manually (no useradd package needed)
RUN mkdir -p /app /tmp && \
    echo 'user:x:10001:10001::/nonexistent:/sbin/nologin' > /etc/passwd && \
    echo 'user:x:10001:' > /etc/group && \
    chmod 1777 /tmp

# binary only
COPY --from=builder /out/app /app/app

# ownership + immutable style perms
RUN chown -R 10001:10001 /app && \
    chmod 0555 /app/app && \
    chmod 0555 /app

# remove shells/tools commonly abused for reverse shell / privilege abuse
RUN rm -f /bin/sh \
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
          /usr/bin/lua \
          /usr/bin/php 2>/dev/null || true

# run as unprivileged user
USER 10001:10001

WORKDIR /app

EXPOSE 3000/tcp

# read-only style app process
ENTRYPOINT ["/app/app"]