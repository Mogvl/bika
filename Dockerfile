# ===========================================
# 多阶段构建: 哔咔漫画 Web 版 (PicACG Web)
# ===========================================

# ---- 阶段1: 构建前端 ----
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm ci

COPY frontend/ .
RUN npm run build:only

# ---- 阶段2: 构建 Go 后端 ----
FROM golang:1.22-alpine AS backend-builder

WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .
# 复制前端构建产物到静态目录
COPY --from=frontend-builder /app/frontend/dist ./static

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bika-server .

# ---- 阶段3: 运行 ----
FROM alpine:3.19

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 创建用户
RUN adduser -D -u 1000 bika

# 创建数据目录
RUN mkdir -p /data/downloads /data/config /data/cache && \
    chown -R bika:bika /data

WORKDIR /app

COPY --from=backend-builder /app/bika-server .
COPY --from=backend-builder /app/static ./static

USER bika

# 环境变量
ENV PORT=8080 \
    DATA_DIR=/data

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:8080/api/health || exit 1

ENTRYPOINT ["/app/bika-server"]
