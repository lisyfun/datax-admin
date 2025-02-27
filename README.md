# DataX Admin

DataX Admin 是一个强大的数据同步和管理平台，基于 Go 和 Vue.js 构建。它提供了友好的用户界面，用于配置、监控和管理 DataX 数据同步任务。

## 功能特性

- 可视化配置 DataX 任务
- 实时监控任务执行状态
- 支持多种数据源（MySQL、PostgreSQL、Oracle 等）
- 任务调度和定时执行
- 用户权限管理
- 系统监控和性能分析
- SSH 终端支持

## 系统要求

- Go 1.18 或更高版本
- Node.js 16 或更高版本
- MySQL 5.7 或更高版本
- DataX 环境

## 快速开始


### 从源码构建

1. 克隆仓库

```bash
git clone https://github.com/lisyfun/datax-admin.git
cd datax-admin
```

2. 使用 Makefile 构建

```bash
# 完整构建（前端 + 后端 + Docker）
make full

# 或者分步构建
make build-frontend  # 构建前端
make build           # 构建后端
```

3. 运行应用

```bash
make run
```

## 使用 Makefile

项目提供了一个功能丰富的 Makefile，用于简化构建和开发过程。

### 主要命令

- `make build` - 构建后端可执行文件（当前平台）
- `make build-frontend` - 构建前端
- `make clean` - 清理构建产物
- `make test` - 运行测试
- `make docker` - 构建 Docker 镜像 (AMD64)
- `make docker-arm64` - 构建 Docker 镜像 (ARM64)
- `make docker-all` - 构建所有架构的 Docker 镜像
- `make docker-multi` - 构建并推送多架构 Docker 镜像
- `make release` - 构建所有平台的发布版本
- `make run` - 构建并运行应用
- `make help` - 显示帮助信息

### 多平台构建

- `make linux-amd64` - 构建 Linux AMD64 可执行文件
- `make linux-arm64` - 构建 Linux ARM64 可执行文件
- `make darwin-amd64` - 构建 macOS AMD64 可执行文件
- `make darwin-arm64` - 构建 macOS ARM64 可执行文件
- `make windows-amd64` - 构建 Windows AMD64 可执行文件

### 示例

```bash
# 构建当前平台的可执行文件
make build

# 构建所有平台的发布版本
make release

# 构建并标记特定版本的 Docker 镜像 (AMD64)
make docker VERSION=v1.0.0

# 构建并标记特定版本的 Docker 镜像 (ARM64)
make docker-arm64 VERSION=v1.0.0

# 构建所有架构的 Docker 镜像
make docker-all VERSION=v1.0.0

# 构建并推送多架构 Docker 镜像
make docker-multi VERSION=v1.0.0

# 完整构建（前端、后端和所有架构的 Docker 镜像）
make full
```

### Docker 镜像说明

项目支持构建不同架构的 Docker 镜像：

- `datax-admin:版本号` - 基于 AMD64 架构的镜像
- `datax-admin:版本号-arm64` - 基于 ARM64 架构的镜像

使用多架构镜像时，Docker 会自动选择适合当前系统的镜像版本。

## 配置

系统配置位于 `backend/config.yaml` 文件中，主要包括：

```yaml
server:
  port: ":28080"
  mode: debug
  base_path: "/"

db:
  host: localhost
  port: 3306
  username: root
  password: password
  dbname: datax_admin
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600
  log_mode: info
  charset: utf8mb4

jwt:
  secret: your-secret-key
  expire: 86400  # 24 hours

datax:
  home: "/app"  # DataX 安装目录
  bin: "/app/bin/datax"  # DataX 可执行文件路径

logger:
  log_path: "logs"  # 日志文件路径
  max_size: 100     # 单个日志文件最大大小（MB）
  max_backups: 7    # 最大保留的旧日志文件数
  max_age: 30       # 日志文件保留的最大天数
  compress: true    # 是否压缩旧日志文件
```

## 项目结构

```
datax-admin/
├── frontend/                # 前端项目目录
│   ├── src/
│   │   ├── api/            # API 接口
│   │   ├── assets/         # 静态资源
│   │   ├── components/     # 公共组件
│   │   ├── layouts/        # 布局组件
│   │   ├── locale/         # 国际化
│   │   ├── router/         # 路由配置
│   │   ├── stores/         # 状态管理
│   │   ├── styles/         # 全局样式
│   │   ├── types/          # TypeScript 类型
│   │   ├── utils/          # 工具函数
│   │   └── views/          # 页面组件
│   ├── package.json
│   └── vite.config.ts
│
└── backend/                 # 后端项目目录
    ├── config/             # 配置文件
    ├── controllers/        # 控制器
    ├── middleware/         # 中间件
    ├── models/             # 数据模型
    ├── routes/             # 路由配置
    ├── services/           # 业务逻辑
    ├── types/              # 类型定义
    ├── utils/              # 工具函数
    ├── go.mod
    └── main.go
```

## 开发指南

### 代码规范

- 前端遵循 Vue 3 + TypeScript 开发规范
- 后端遵循 Go 语言规范
- 使用 ESLint 和 Prettier 进行代码格式化
- 使用 golangci-lint 进行 Go 代码检查

### 提交规范

提交信息格式：
```
<type>(<scope>): <subject>

<body>

<footer>
```

type 类型：
- feat: 新功能
- fix: 修复
- docs: 文档
- style: 格式
- refactor: 重构
- test: 测试
- chore: 构建过程或辅助工具的变动

### 分支管理

- main: 主分支，用于生产环境
- develop: 开发分支
- feature/*: 功能分支
- hotfix/*: 紧急修复分支

## 部署指南

### Docker 部署

1. 构建镜像
```bash
# 前端
docker build -t datax-admin-frontend ./frontend

# 后端
docker build -t datax-admin-backend ./backend
```

2. 运行容器
```bash
# 前端
docker run -d -p 80:80 datax-admin-frontend

# 后端
docker run -d -p 8080:8080 datax-admin-backend
```

### 传统部署

1. 前端部署
   - 执行 `pnpm build` 生成静态文件
   - 将 `dist` 目录下的文件部署到 Web 服务器

2. 后端部署
   - 执行 `go build` 生成可执行文件
   - 配置系统服务或使用 PM2 等工具运行

## 贡献指南

1. Fork 本仓库
2. 创建功能分支
3. 提交代码
4. 创建 Pull Request

## 许可证

[MIT License](LICENSE)
