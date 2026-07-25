FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata wget

RUN adduser -D -u 1000 bika && \
    mkdir -p /data/downloads /data/config /data/cache && \
    chown -R bika:bika /data

WORKDIR /app

COPY backend/bika-server .
COPY backend/static ./static

USER bika

ENV PORT=4000 \
    DATA_DIR=/data

EXPOSE 4000

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:${PORT:-4000}/api/health || exit 1

ENTRYPOINT ["/app/bika-server"]
