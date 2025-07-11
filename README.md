# DataX Admin



![image](https://github.com/user-attachments/assets/8a07b6e3-7c47-4f23-9fd9-ff7c7d802787)

DataX Admin 是一个基于 Go 和 Vue.js 构建的现代化数据同步任务管理平台，提供了友好的界面来配置和管理 DataX 数据同步任务。

## 功能特性

### 核心功能
- 🎯 **可视化配置 DataX 任务** - 支持多种数据源的可视化配置
- 📊 **实时监控任务执行状态** - 任务执行状态实时监控和历史记录
- 🗄️ **多数据源支持** - 支持 MySQL、PostgreSQL、Oracle 等
- ⏰ **任务调度和定时执行** - 基于 Cron 表达式的灵活任务调度
- 🔐 **完整的权限管理系统** - 基于角色的权限控制，支持菜单和按钮级别权限
- 🖥️ **终端管理** - SSH 终端连接管理，支持密码和密钥文件认证
- 📈 **系统监控和性能分析** - 实时系统资源监控和性能分析
- 🐳 **容器化部署** - 支持多架构部署（AMD64/ARM64）

### 技术特性
- 🎨 **现代化 UI** - 基于 Arco Design Vue 的现代化界面
- 🔒 **安全加密** - 支持数据传输加密和安全认证
- 🚀 **热重载开发** - 开发环境支持前后端热重载
- 📱 **响应式设计** - 适配各种屏幕尺寸
- 🔧 **环境变量配置** - 支持 Docker 环境变量配置

## 系统要求

- Go 1.20 或更高版本
- Node.js 18 或更高版本
- 数据库（二选一）：
  - MySQL 8.0 或更高版本
  - PostgreSQL 12 或更高版本
- DataX 环境

## 快速开始

### 使用 Docker（推荐）

#### 方式一：Docker Compose（最简单）

**使用 PostgreSQL 数据库：**
```bash
# 下载 PostgreSQL 版本的 docker-compose 文件
curl -O https://raw.githubusercontent.com/lisyfun/datax-admin/main/docker-compose-postgres.yml

# 启动服务（包含 PostgreSQL 数据库）
docker-compose -f docker-compose-postgres.yml up -d
```

**使用 MySQL 数据库：**
```bash
# 下载 MySQL 版本的 docker-compose 文件
curl -O https://raw.githubusercontent.com/lisyfun/datax-admin/main/docker-compose-mysql.yml

# 启动服务（包含 MySQL 数据库）
docker-compose -f docker-compose-mysql.yml up -d
```

#### 方式二：单独运行容器

1. 拉取镜像
```bash
docker pull lisongyu/datax-admin:latest
```

2. 运行容器

**使用 MySQL 数据库：**
```bash
docker run -d \
  --name datax-admin \
  -p 28080:80 \
  -v $(pwd)/logs:/app/logs \
  -e DB_TYPE=mysql \
  -e DB_HOST=localhost \
  -e DB_PORT=3306 \
  -e DB_USERNAME=root \
  -e DB_PASSWORD=your_password \
  -e DB_NAME=datax_admin \
  -e JWT_SECRET=your-jwt-secret-key \
  lisongyu/datax-admin:latest
```

**使用 PostgreSQL 数据库：**
```bash
docker run -d \
  --name datax-admin \
  -p 28080:80 \
  -v $(pwd)/logs:/app/logs \
  -e DB_TYPE=postgres \
  -e DB_HOST=localhost \
  -e DB_PORT=5432 \
  -e DB_USERNAME=postgres \
  -e DB_PASSWORD=your_password \
  -e DB_NAME=datax_admin \
  -e DB_SSL_MODE=disable \
  -e DB_TIME_ZONE=Asia/Shanghai \
  -e JWT_SECRET=your-jwt-secret-key \
  lisongyu/datax-admin:latest
```

**环境变量说明：**
- `DB_TYPE`: 数据库类型（默认：mysql，支持：mysql, postgres）
- `DB_HOST`: 数据库主机地址（默认：localhost）
- `DB_PORT`: 数据库端口（MySQL默认：3306，PostgreSQL默认：5432）
- `DB_USERNAME`: 数据库用户名（默认：root）
- `DB_PASSWORD`: 数据库密码（必须设置）
- `DB_NAME`: 数据库名称（默认：datax_admin）
- `DB_SSL_MODE`: PostgreSQL SSL模式（默认：disable）
- `DB_TIME_ZONE`: 数据库时区（默认：Asia/Shanghai）
- `JWT_SECRET`: JWT密钥（建议设置复杂密钥）

**注意：** 请确保数据库已创建并可访问，系统会自动创建所需的表结构。

### 数据库初始化

系统支持 **MySQL** 和 **PostgreSQL** 两种数据库，在运行系统之前，请确保：

#### MySQL 数据库
1. MySQL 数据库服务已启动
2. 创建数据库（如果不存在）：
```sql
CREATE DATABASE datax_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```
3. 确保数据库用户有足够的权限访问该数据库
4. 在 `config.yaml` 中设置 `db.type: mysql`

#### PostgreSQL 数据库
1. PostgreSQL 数据库服务已启动
2. 创建数据库（如果不存在）：
```sql
CREATE DATABASE datax_admin;
```
3. 确保数据库用户有足够的权限访问该数据库
4. 在 `config.yaml` 中设置 `db.type: postgres`
5. 可参考 `config-postgres.yaml` 配置示例

系统首次启动时会自动创建所需的表结构和初始数据。

#### 强制重新初始化

如果需要重新初始化数据（会清空现有数据），可以在配置文件中设置：

```yaml
init:
  force_reset: true
```

**⚠️ 警告：** 此操作会删除所有现有的用户、角色和权限数据，请谨慎使用！

详细说明请参考：[强制重置配置说明](FORCE_RESET_CONFIG.md)

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
  -e DB_HOST=localhost \
  -e DB_PORT=3306 \
  -e DB_USERNAME=root \
  -e DB_PASSWORD=your_password \
  -e DB_NAME=datax_admin \
  -e JWT_SECRET=your-jwt-secret-key \
  datax-admin:$(VERSION)-arm64
```


## 开发和构建

项目提供了完整的 Makefile 来简化开发和构建过程。

### 开发环境

**启动开发服务（推荐）**
```bash
make run
```
此命令会同时启动前端和后端服务，支持热重载：
- 后端服务：http://localhost:28080
- 前端服务：http://localhost:3000

### 构建命令

- `make build` - 构建后端（当前平台）
- `make build-arm64` - 构建后端（ARM64）
- `make build-frontend` - 构建前端
- `make docker VERSION=v1.0.0` - 构建 Docker 镜像（AMD64）
- `make docker-arm64 VERSION=v1.0.0` - 构建 Docker 镜像（ARM64）
- `make docker-all VERSION=v1.0.0` - 构建所有架构的 Docker 镜像
- `make clean` - 清理构建产物
- `make help` - 显示帮助信息


## 配置说明

### 环境变量配置（Docker 推荐）

使用 Docker 运行时，可以通过环境变量进行配置：

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `DB_TYPE` | 数据库类型 | mysql |
| `DB_HOST` | 数据库主机地址 | localhost |
| `DB_PORT` | 数据库端口 | 3306 (MySQL) / 5432 (PostgreSQL) |
| `DB_USERNAME` | 数据库用户名 | root |
| `DB_PASSWORD` | 数据库密码 | 无（必须设置） |
| `DB_NAME` | 数据库名称 | datax_admin |
| `DB_SSL_MODE` | PostgreSQL SSL模式 | disable |
| `DB_TIME_ZONE` | 数据库时区 | Asia/Shanghai |
| `JWT_SECRET` | JWT密钥 | 无（建议设置） |
| `SERVER_PORT` | 服务端口 | 28080 |
| `DATAX_HOME` | DataX-Admin 可执行文件路径 | /app |
| `DATAX_BIN` | DataX 可执行文件路径 | /app/bin/datax |
| `DATAX_PYTHON` | Python 解释器路径 | python |

### 配置文件（源码部署）

系统配置文件位于 `config.yaml`：

```yaml
server:
  port: ":28080"
  mode: release
  base_path: "/"
  max_file_size: 500  # 500MB

db:
  type: mysql  # 数据库类型: mysql, postgres
  host: localhost
  port: 3306   # MySQL: 3306, PostgreSQL: 5432
  username: root
  password: your_password
  dbname: datax_admin
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600
  log_mode: warn
  charset: utf8mb4        # MySQL专用
  time_zone: Asia/Shanghai
  ssl_mode: disable       # PostgreSQL专用

datax:
  home: "/app"  # DataX 安装目录
  bin: "/app/bin/datax"  # DataX 可执行文件路径

# 认证配置
auth:
  secret: "your-auth-secret-key"    # 认证密钥
  expiration: 86400                 # 过期时间(秒)，86400=1天

logger:
  log_path: "logs"  # 日志文件路径
  max_size: 100     # 单个日志文件最大大小（MB）
  max_backups: 7    # 最大保留的旧日志文件数
  max_age: 30       # 日志文件保留的最大天数
  compress: true    # 是否压缩旧日志文件

# 初始化配置
init:
  force_reset: false # 是否强制重新初始化数据（会清空现有数据）
```

## 项目结构

```
datax-admin/
├── web/                     # 前端项目目录 (Vue 3 + TypeScript + Arco Design)
│   ├── src/
│   │   ├── api/            # API 接口定义
│   │   ├── assets/         # 静态资源
│   │   ├── components/     # 公共组件
│   │   ├── directives/     # 自定义指令（权限指令等）
│   │   ├── layouts/        # 布局组件
│   │   ├── router/         # 路由配置
│   │   ├── stores/         # 状态管理 (Pinia)
│   │   ├── styles/         # 全局样式
│   │   ├── types/          # TypeScript 类型定义
│   │   ├── utils/          # 工具函数
│   │   └── views/          # 页面组件
│   │       ├── dashboard/  # 仪表盘
│   │       ├── job/        # 任务管理
│   │       ├── terminal/   # 终端管理
│   │       ├── system/     # 系统管理
│   │       └── login/      # 登录页面
│   ├── package.json
│   ├── vite.config.ts
│   └── index.html
│
├── config/                 # 配置管理
├── controllers/            # 控制器层
├── datax/                  # DataX 相关文件
├── docker/                 # Docker 配置文件
├── middleware/             # 中间件（认证、安全、日志等）
├── models/                 # 数据模型 (GORM)
├── routes/                 # 路由配置
├── services/               # 业务逻辑层
├── types/                  # Go 类型定义
├── utils/                  # 工具函数
├── bin/                    # 编译后的可执行文件
├── logs/                   # 日志文件目录
├── Makefile               # 构建脚本
├── go.mod                 # Go 模块文件
├── go.sum                 # Go 模块校验文件
├── main.go                # 主入口文件
└── config.yaml            # 配置文件
```

## 开发指南

### 技术栈

**前端技术栈：**
- Vue 3 + TypeScript
- Arco Design Vue（UI 组件库）
- Vite（构建工具）
- Pinia（状态管理）
- Vue Router（路由管理）
- Axios（HTTP 客户端）

**后端技术栈：**
- Go 1.20+
- Gin（Web 框架）
- GORM（ORM 框架）
- MySQL（数据库）
- Cron（任务调度）

### 开发环境设置

1. **克隆项目**
```bash
git clone https://github.com/lisyfun/datax-admin.git
cd datax-admin
```

2. **启动开发服务**
```bash
make run
```

3. **访问应用**
- 前端开发服务：http://localhost:3000/datax/
- 后端 API 服务：http://localhost:28080

### 代码规范

- 前端遵循 Vue 3 + TypeScript 开发规范
- 后端遵循 Go 语言规范
- 使用 ESLint 和 Prettier 进行代码格式化
- 使用 golangci-lint 进行 Go 代码检查

### 主要功能模块

- **仪表盘**：系统概览、统计信息、系统监控
- **任务管理**：DataX 任务的创建、编辑、执行、监控
- **终端管理**：SSH 终端连接管理，支持密码和密钥认证
- **系统管理**：用户管理、角色管理、权限管理
- **权限系统**：基于角色的访问控制，支持菜单和按钮级别权限


## 贡献指南

1. Fork 本仓库
2. 创建功能分支
3. 提交代码
4. 创建 Pull Request

## 许可证

[MIT License](LICENSE)
