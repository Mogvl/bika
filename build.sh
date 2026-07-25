#!/bin/bash
# 哔咔漫画 Web 版 - Docker 构建脚本
# 用法: ./build.sh [架构]
# 架构: amd64 (默认, Intel/AMD), arm64 (ARM, 如部分绿联云 NAS)

set -e

ARCH=${1:-amd64}
IMAGE_NAME="bika-web:latest"

echo "=========================================="
echo "  构建哔咔漫画 Docker 镜像"
echo "  架构: ${ARCH}"
echo "  镜像: ${IMAGE_NAME}"
echo "=========================================="

# 检查 docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "错误: 未找到 docker 命令，请先安装 Docker"
    exit 1
fi

# 构建镜像
echo ""
echo ">>> 开始构建（首次构建约需 3-5 分钟）..."
docker build \
    --build-arg TARGETARCH=${ARCH} \
    -t ${IMAGE_NAME} \
    -f Dockerfile \
    .

echo ""
echo "=========================================="
echo "  ✅ 构建成功!"
echo ""
echo "  部署方式:"
echo "  1. 打开绿联云 Docker → Compose"
echo "  2. 上传 docker-compose.yml"
echo "  3. 点击部署"
echo ""
echo "  或命令行启动:"
echo "  docker compose up -d"
echo "=========================================="
