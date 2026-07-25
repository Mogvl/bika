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

```bash
# 克隆仓库
git clone https://github.com/Mogvl/bika.git
cd bika

# 启动服务
docker compose up -d

# 查看日志
docker compose logs -f
```

或者直接新建 `docker-compose.yml` 文件，粘贴以下内容：

```yaml
version: '3.8'

services:
  bika:
    build: .
    container_name: bika-web
    restart: unless-stopped

    # ==================== 端口映射 ====================
    # 格式: "宿主机端口:容器端口"
    # 宿主机端口可自定义，避免与其它服务冲突
    # 绿联云访问地址: http://<NAS-IP>:4000
    ports:
      - "4000:4000"

    # ==================== 环境变量 ====================
    # PORT: 服务端口（需与上方容器端口一致）
    # TZ: 时区设置
    environment:
      - PORT=4000
      - TZ=Asia/Shanghai

    # ==================== 数据持久化 ====================
    # 将容器内数据映射到宿主机，防止重启丢失
    # ./data 目录会自动创建在当前文件夹下
    volumes:
      # 漫画下载保存位置
      - ./data/downloads:/data/downloads
      # 配置文件持久化
      - ./data/config:/data/config
      # 缓存文件
      - ./data/cache:/data/cache

    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:4000/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
```

然后执行：

```bash
docker compose up -d
```

访问 `http://localhost:4000` 即可使用。

### 使用 Docker 命令行

```bash
docker build -t bika-web .
docker run -d \
  --name bika-web \
  -p 4000:4000 \
  -e PORT=4000 \
  -e TZ=Asia/Shanghai \
  -v ./data/downloads:/data/downloads \
  -v ./data/config:/data/config \
  -v ./data/cache:/data/cache \
  --restart unless-stopped \
  bika-web
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

1. 打开绿联云 NAS 的 **文件管理器**
2. 创建一个文件夹，例如 `docker/bika`
3. 将本项目中的 `docker-compose.yml` 上传到该目录
4. 打开 **Docker** 应用
5. 进入 **Compose** 页面
6. 选择项目目录为 `docker/bika`
7. 点击 **部署**

### 方法2：Portainer 部署

1. 在绿联云上安装 Portainer
2. 进入 Portainer 管理界面
3. 选择 **Stacks** → **Add stack**
4. 粘贴 `docker-compose.yml` 内容
5. 点击 **Deploy**

### 方法3：命令行部署

通过 SSH 连接绿联云 NAS，执行：

```bash
# 拉取代码
git clone https://github.com/Mogvl/bika.git /volume1/docker/bika
cd /volume1/docker/bika

# 构建并启动
docker compose up -d
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
