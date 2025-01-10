# DataX Admin

DataX Admin 是一个现代化的前后端分离管理系统，采用 Vue 3 + Go 技术栈开发。

## 技术栈

### 前端技术栈

- Vue 3：渐进式 JavaScript 框架
- TypeScript：JavaScript 的超集，提供类型检查
- Arco Design Vue：字节跳动出品的企业级设计系统
- Pinia：Vue 的状态管理库
- Vue Router：Vue.js 的官方路由
- Vite：新一代前端构建工具
- i18n：国际化解决方案

### 后端技术栈

- Go 1.22：编程语言
- Gin：Web 框架
- GORM：ORM 框架
- JWT：用户认证
- MySQL：数据库
- Viper：配置管理

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

## 功能特性

- 用户认证与授权
- 菜单管理
- 角色管理
- 权限管理
- 系统配置
- 国际化支持
- 主题定制
- 响应式设计

## 快速开始

### 环境要求

- Node.js >= 16
- Go >= 1.22
- MySQL >= 8.0

### 前端启动

```bash
# 进入前端目录
cd frontend

# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev

# 构建生产版本
pnpm build
```

### 后端启动

```bash
# 进入后端目录
cd backend

# 安装依赖
go mod tidy

# 启动服务
go run main.go
```

### 配置说明

1. 前端配置
   - 环境变量配置在 `.env` 文件中
   - 开发环境配置在 `.env.development` 文件中
   - 生产环境配置在 `.env.production` 文件中

2. 后端配置
   - 配置文件位于 `backend/config` 目录
   - 支持 YAML 格式的配置文件
   - 可通过环境变量覆盖配置

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
