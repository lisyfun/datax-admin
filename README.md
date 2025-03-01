# DataX Admin



![image](https://github.com/user-attachments/assets/7882709f-6147-47cc-ab30-4f4515538ff4)


DataX Admin 是一个基于 Go 和 Vue.js 构建的数据同步任务管理平台，提供了友好的界面来配置和管理 DataX 数据同步任务。

## 功能特性

- 可视化配置 DataX 任务
- 实时监控任务执行状态
- 支持多种数据源管理
- 任务调度和定时执行
- 用户权限管理
- 系统监控和性能分析
- 支持多架构部署（AMD64/ARM64）

## 系统要求

- Go 1.20 或更高版本
- Node.js 18 或更高版本
- MySQL 8.0 或更高版本
- DataX 环境

## 快速开始

### 使用 Docker（推荐）

1. 拉取镜像
```bash
docker pull lisongyu/datax-admin:latest
```

2. 运行容器
```bash
docker run -d \
  --name datax-admin \
  -p 28080:80 \
  -v $(pwd)/logs:/app/logs \
  lisongyu/datax-admin:latest  # ARM64 使用 datax-admin:latest-arm64
```

访问 http://localhost:28080/datax 即可使用系统。

### 从源码构建

1. 克隆仓库
```bash
git clone https://github.com/lisyfun/datax-admin.git
cd datax-admin
```

2. 使用 Makefile 构建（推荐）
```bash
# 完整构建（前端 + 后端 + Docker）
make docker
```

3. 运行容器
```bash
docker run -d \
  --name datax-admin \
  -p 28080:80 \
  -v $(pwd)/logs:/app/logs \
  lisongyu/datax-admin:latest  # ARM64 使用 datax-admin:latest-arm64
```


## 使用 Makefile

项目提供了完整的 Makefile 来简化构建过程。

### 主要命令

- `make build` - 构建后端（当前平台）
- `make build-frontend` - 构建前端
- `make docker VERSION=v1.0.0` - 构建 Docker 镜像
- `make docker-run VERSION=v1.0.0` - 运行 Docker 容器
- `make clean` - 清理构建产物
- `make help` - 显示帮助信息


## 配置说明

系统配置文件位于 `backend/config.yaml`：

```yaml
server:
  port: ":28080"
  mode: debug

db:
  host: localhost
  port: 3306
  username: root
  password: password
  dbname: datax_admin

jwt:
  secret: your-secret-key
  expire: 86400

datax:
  home: "/app"  # DataX-Admin 可执行文件路径
  bin: "/app/bin/datax"  # DataX 可执行文件路径
  python: "python"  # Python 解释器路径

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
