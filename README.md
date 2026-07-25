# 哔咔漫画 Web 版 (PicACG Web)

基于 [tonquer/picacg-qt](https://github.com/tonquer/picacg-qt) 的 Web 版本，使用 **Go + Vue 3** 重构，支持 Docker 部署，专为 **绿联云 NAS** 等设备优化。

## 功能

- ✅ 哔咔漫画账号登录
- ✅ 分类浏览漫画
- ✅ 搜索漫画（支持热词）
- ✅ 漫画详情查看
- ✅ 在线漫画阅读器（单页/滚动模式）
- ✅ 排行榜（24小时/本周/本月）
- ✅ 收藏管理
- ✅ 图片代理（解决跨域和缓存）
- ✅ Docker 多阶段构建（镜像小巧）

## 快速开始

### 使用 Docker Compose（推荐）

> **注意**: 绿联云 Docker Compose 不支持 `build` 指令，需要先手动构建镜像。

**第一步：构建镜像**

通过 SSH 连接绿联云 NAS，或在项目目录下执行：

```bash
# Intel/AMD 架构（大部分绿联云）
./build.sh amd64

# ARM 架构（部分绿联云型号）
./build.sh arm64
```

或者直接用 docker 命令：

```bash
docker build -t bika-web:latest .
```

**第二步：部署**

新建 `docker-compose.yml` 文件，粘贴以下内容：

```yaml
version: '3.8'

services:
  bika:
    image: bika-web:latest
    container_name: bika-web
    restart: unless-stopped
    ports:
      - "4000:4000"
    environment:
      - PORT=4000
      - TZ=Asia/Shanghai
    volumes:
      - ./data/downloads:/data/downloads
      - ./data/config:/data/config
      - ./data/cache:/data/cache
```

然后启动：

```bash
docker compose up -d
```

访问 `http://localhost:4000` 即可使用。

### 使用 Docker 命令行（无需 Compose）

```bash
# 先构建镜像
docker build -t bika-web:latest .

# 运行容器
docker run -d \
  --name bika-web \
  --restart unless-stopped \
  -p 4000:4000 \
  -e PORT=4000 \
  -e TZ=Asia/Shanghai \
  -v ./data/downloads:/data/downloads \
  -v ./data/config:/data/config \
  -v ./data/cache:/data/cache \
  bika-web:latest
```

### 本地开发

```bash
# 启动 Go 后端
cd backend
go run .

# 新终端，启动 Vue 前端开发服务器
cd frontend
npm install
npm run dev
```

前端开发服务器默认监听 `http://localhost:3000`，API 请求会代理到 `http://localhost:4000`。

## 绿联云 NAS 部署指南

### 方法1：Docker Compose（推荐）

1. **通过 SSH 连接 NAS**（或在 Docker 终端执行）
2. 拉取代码并构建镜像：
   ```bash
   git clone https://github.com/Mogvl/bika.git /volume1/docker/bika
   cd /volume1/docker/bika
   docker build -t bika-web:latest .
   ```
3. 打开绿联云 NAS 的 **Docker** 应用
4. 进入 **Compose** 页面
5. 新建项目，粘贴以下内容（项目目录选 `/volume1/docker/bika`）：
   ```yaml
   version: '3.8'
   services:
     bika:
       image: bika-web:latest
       container_name: bika-web
       restart: unless-stopped
       ports:
         - "4000:4000"
       environment:
         - PORT=4000
         - TZ=Asia/Shanghai
       volumes:
         - /volume1/docker/bika/data/downloads:/data/downloads
         - /volume1/docker/bika/data/config:/data/config
         - /volume1/docker/bika/data/cache:/data/cache
   ```
6. 点击 **部署**

### 方法2：Portainer 部署

1. 在绿联云上安装 Portainer
2. 进入 Portainer 管理界面
3. 选择 **Images** → **Build a new image**，上传项目目录构建
4. 选择 **Stacks** → **Add stack**
5. 粘贴 `docker-compose.yml` 内容（将 `build: .` 改为 `image: bika-web:latest`）
6. 点击 **Deploy**

### 方法3：命令行部署（最简单）

通过 SSH 连接绿联云 NAS，直接执行：

```bash
# 拉取代码
git clone https://github.com/Mogvl/bika.git /volume1/docker/bika
cd /volume1/docker/bika

# 构建镜像
docker build -t bika-web:latest .

# 创建数据目录
mkdir -p data/{downloads,config,cache}

# 运行容器
docker run -d \
  --name bika-web \
  --restart unless-stopped \
  -p 4000:4000 \
  -e PORT=4000 \
  -e TZ=Asia/Shanghai \
  -v /volume1/docker/bika/data/downloads:/data/downloads \
  -v /volume1/docker/bika/data/config:/data/config \
  -v /volume1/docker/bika/data/cache:/data/cache \
  bika-web:latest
```

### 访问地址

部署成功后，通过以下地址访问：

```
http://<NAS-IP>:4000
```

默认端口为 `4000`，可在 `docker-compose.yml` 中修改映射端口。

## 配置说明

### 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | `8080` | 服务监听端口 |
| `TZ` | `Asia/Shanghai` | 时区设置 |

### 数据持久化

通过 Docker 卷挂载持久化数据：

| 挂载路径 | 说明 |
|----------|------|
| `./data/downloads` | 漫画下载目录 |
| `./data/config` | 配置保存目录 |
| `./data/cache` | 缓存目录 |

## 登录说明

使用 **哔咔漫画** 的邮箱和密码登录。

> 首次使用需要注册哔咔漫画账号（通过哔咔漫画 App 注册）。

## 技术栈

- **后端**: Go (标准库 `net/http`)
- **前端**: Vue 3 + TypeScript + Vite
- **状态管理**: Pinia
- **路由**: Vue Router
- **HTTP 客户端**: Axios
- **容器**: Docker 多阶段构建

## 与原始项目的关系

本项目是 [tonquer/picacg-qt](https://github.com/tonquer/picacg-qt)（PySide6 桌面客户端）的 Web 重构版本：

- 保留了所有核心的哔咔漫画 API 交互逻辑
- 将桌面 GUI 重构为 Web 界面
- 使用 Go 重写了后端 API 代理层
- 支持 Docker 部署，适合 NAS 环境
- 保留图片代理功能以支持浏览器访问

## 免责声明

- 本项目仅供技术研究使用
- 请勿用于商业用途
- 使用哔咔漫画 API 需遵守平台相关规定
- 如有侵权，请联系删除

## License

LGPL-3.0
