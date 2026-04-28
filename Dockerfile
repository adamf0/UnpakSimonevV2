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
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

RUN microdnf update -y && \
    microdnf install -y tzdata ca-certificates shadow-utils && \
    microdnf clean all

ENV TZ=Asia/Jakarta \
    HOME=/nonexistent \
    PATH=/app

RUN mkdir -p /app /tmp && \
    useradd -r -u 10001 -g 0 user && \
    chmod 1777 /tmp

COPY --from=builder /src/out/app /app/app

RUN chown -R 10001:0 /app && \
    chmod 0555 /app/app && \
    chmod -R g=u /app

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
    
USER 10001

WORKDIR /app

EXPOSE 3000

ENTRYPOINT ["/app/app"]