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

## 绿联云 NAS 部署（推荐）

只需两步，通过 Compose 自动拉取镜像部署：

### 第一步：在 GitHub Actions 构建镜像

1. 打开 [GitHub Actions](https://github.com/Mogvl/bika/actions)
2. 选择 **构建 Docker 镜像** 工作流
3. 点击 **Run workflow** 手动触发构建
4. 等待构建完成（约 3-5 分钟）

### 第二步：在绿联云部署 Compose

1. 打开绿联云 NAS 的 **Docker** 应用
2. 进入 **Compose** 页面
3. 新建项目，粘贴以下内容：

```yaml
services:
  bika-web:
    image: ghcr.io/mogvl/bika-web:latest
    container_name: bika-web
    ports:
      - "3000:4000"
    environment:
      - PORT=4000
      - TZ=Asia/Shanghai
    volumes:
      - /volume1/bika:/data/downloads
      - /volume1/docker/bika:/data/config
      - /volume1/docker/bika/cache:/data/cache
    restart: unless-stopped
```

4. 点击 **部署**

部署成功后访问 `http://<NAS-IP>:3000` 即可使用。

## 快速开始（本地开发）

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

## 使用 Docker 命令行（无需 Compose）

```bash
docker run -d \
  --name bika-web \
  --restart unless-stopped \
  -p 3000:4000 \
  -e PORT=4000 \
  -e TZ=Asia/Shanghai \
  -v /volume1/bika:/data/downloads \
  -v /volume1/docker/bika:/data/config \
  -v /volume1/docker/bika/cache:/data/cache \
  ghcr.io/mogvl/bika-web:latest
```

## 配置说明

### 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | `4000` | 服务监听端口 |
| `TZ` | `Asia/Shanghai` | 时区设置 |

### 数据持久化

| 挂载路径 | 说明 |
|----------|------|
| `/volume1/bika` | 漫画下载目录 |
| `/volume1/docker/bika` | 配置保存目录 |
| `/volume1/docker/bika/cache` | 缓存目录 |

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
- **CI/CD**: GitHub Actions

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
