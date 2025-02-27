# 第一阶段：最终镜像
FROM caddy:2.7-alpine

# 设置工作目录
WORKDIR /usr/share/caddy

# 环境变量
ENV TZ=Asia/Shanghai \
    LANG=C.UTF-8

# 安装依赖
RUN apk add --no-cache tzdata ca-certificates bash curl && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 参数定义
ARG BINARY_NAME=datax-admin
ARG BINARY_VERSION=linux-amd64
ARG DATAX_VERSION=linux-amd64
ARG CONFIG_FILE=./backend/config.yaml

# 创建必要的目录
RUN mkdir -p /app/logs /app/bin

# 复制前端构建产物
COPY ./frontend/dist ./dist/

# 复制后端二进制文件和配置
COPY ./backend/bin/${BINARY_NAME}-${BINARY_VERSION} /app/${BINARY_NAME}
COPY ./bin/datax-${DATAX_VERSION} /app/bin/datax
COPY ${CONFIG_FILE} /app/config.yaml

# 复制 Caddy 配置和启动脚本
COPY Caddyfile /etc/caddy/Caddyfile
COPY start.sh /app/start.sh

# 设置工作目录
WORKDIR /app

# 设置可执行权限
RUN chmod +x /app/${BINARY_NAME} /app/bin/datax /app/start.sh && \
    caddy fmt --overwrite /etc/caddy/Caddyfile

# 暴露端口
EXPOSE 80

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost/datax/api/v1/ping || exit 1

# 启动命令
CMD ["sh", "/app/start.sh"]
