FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata wget

# 创建数据目录（用 root 确保权限正确）
RUN mkdir -p /data/downloads /data/config /data/cache && \
    chmod -R 777 /data

WORKDIR /app

COPY backend/bika-server .
COPY backend/static ./static

ENV PORT=4000 \
    DATA_DIR=/data \
    DOWNLOAD_DIR=/data/downloads

EXPOSE 4000

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:${PORT:-4000}/api/health || exit 1

ENTRYPOINT ["/app/bika-server"]
