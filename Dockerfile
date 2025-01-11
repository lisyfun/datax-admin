# 使用多阶段构建
# 前端构建阶段
FROM node:18-alpine as frontend-builder

# 设置工作目录
WORKDIR /app/frontend

# 复制前端项目文件
COPY frontend/package.json frontend/pnpm-lock.yaml ./

# 安装 pnpm
RUN npm install -g pnpm

# 安装依赖
RUN pnpm install

# 复制前端源代码
COPY frontend/ .

# 构建前端项目（跳过类型检查）
RUN pnpm vite build

# 后端构建阶段
FROM golang:1.22-alpine as backend-builder

# 设置工作目录
WORKDIR /app/backend

# 设置 Go 环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 复制后端项目文件
COPY backend/go.mod backend/go.sum ./

# 下载依赖
RUN go mod download

# 复制后端源代码
COPY backend/ .

# 构建后端项目
RUN go build -o datax-admin .

# 最终运行阶段
FROM alpine:latest

# 安装必要的系统依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 从构建阶段复制构建产物
COPY --from=backend-builder /app/backend/datax-admin .
COPY --from=frontend-builder /app/frontend/dist ./dist
COPY backend/config ./config

# 设置环境变量
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 3000

# 启动命令
CMD ["./datax-admin"]
