# 构建阶段
FROM golang:1.22-alpine as builder

# 设置工作目录
WORKDIR /app

# 设置 Go 环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 复制项目文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建项目
RUN go build -o datax-admin .

# 运行阶段
FROM alpine:latest

# 安装必要的系统依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 从构建阶段复制构建产物
COPY --from=builder /app/datax-admin .
COPY config ./config

# 设置环境变量
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 3000

# 启动命令
CMD ["./datax-admin"]
