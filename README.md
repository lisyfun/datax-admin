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
  -e DB_HOST=localhost \
  -e DB_PORT=3306 \
  -e DB_USERNAME=root \
  -e DB_PASSWORD=your_password \
  -e DB_NAME=datax_admin \
  -e JWT_SECRET=your-jwt-secret-key \
  lisongyu/datax-admin:latest
```

**环境变量说明：**
- `DB_HOST`: 数据库主机地址（默认：localhost）
- `DB_PORT`: 数据库端口（默认：3306）
- `DB_USERNAME`: 数据库用户名（默认：root）
- `DB_PASSWORD`: 数据库密码（必须设置）
- `DB_NAME`: 数据库名称（默认：datax_admin）
- `JWT_SECRET`: JWT密钥（建议设置复杂密钥）

**注意：** 请确保数据库已创建并可访问，系统会自动创建所需的表结构。

### 数据库初始化

在运行系统之前，请确保：

1. MySQL 数据库服务已启动
2. 创建数据库（如果不存在）：
```sql
CREATE DATABASE datax_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```
3. 确保数据库用户有足够的权限访问该数据库

系统首次启动时会自动创建所需的表结构和初始数据。

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


## 使用 Makefile

项目提供了完整的 Makefile 来简化构建过程。

### 主要命令

- `make build` - 构建后端（当前平台）
- `make build-frontend` - 构建前端
- `make docker VERSION=v1.0.0` - 构建 Docker 镜像
- `make docker-run VERSION=v1.0.0` - 运行 Docker 容器（需要配置环境变量）
- `make clean` - 清理构建产物
- `make help` - 显示帮助信息


## 配置说明

### 环境变量配置（Docker 推荐）

使用 Docker 运行时，可以通过环境变量进行配置：

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `DB_HOST` | 数据库主机地址 | localhost |
| `DB_PORT` | 数据库端口 | 3306 |
| `DB_USERNAME` | 数据库用户名 | root |
| `DB_PASSWORD` | 数据库密码 | 无（必须设置） |
| `DB_NAME` | 数据库名称 | datax_admin |
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
├── web/                     # 前端项目目录
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
├── config/                 # 配置文件
├── controllers/            # 控制器
├── datax/                  # DataX相关文件
├── docker/                 # Docker相关文件
├── middleware/             # 中间件
├── models/                 # 数据模型
├── routes/                 # 路由配置
├── services/               # 业务逻辑
├── types/                  # 类型定义
├── utils/                  # 工具函数
├── bin/                    # 可执行文件
├── logs/                   # 日志文件
├── rules/                  # 规则文件
├── go.mod                  # Go 模块文件
├── go.sum                  # Go 模块校验文件
├── main.go                 # 主入口文件
└── config.yaml             # 配置文件
```

## 开发指南

### 代码规范

- 前端遵循 Vue 3 + TypeScript 开发规范
- 后端遵循 Go 语言规范
- 使用 ESLint 和 Prettier 进行代码格式化
- 使用 golangci-lint 进行 Go 代码检查


## 贡献指南

1. Fork 本仓库
2. 创建功能分支
3. 提交代码
4. 创建 Pull Request

## 许可证

[MIT License](LICENSE)
