# 第一阶段：最终镜像
FROM nginx:1.26-alpine

# 设置工作目录
WORKDIR /usr/share/nginx/html

# 环境变量
ENV TZ=Asia/Shanghai \
    LANG=C.UTF-8

# 安装依赖
RUN apk add --no-cache tzdata bash curl sshpass openssh-client && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    mkdir -p /app/logs /app/bin

# 参数定义
ARG BINARY_NAME=datax-admin
ARG BINARY_VERSION=linux-amd64
ARG DATAX_VERSION=linux-amd64
ARG CONFIG_FILE=./config.yaml


# 复制前端构建产物
COPY ./web/dist ./
# 复制错误页面
COPY ./web/404.html ./404.html
COPY ./web/403.html ./403.html
COPY ./web/50x.html ./50x.html

# 复制后端二进制文件和配置
COPY ./bin/${BINARY_NAME}-${BINARY_VERSION} /app/${BINARY_NAME}
COPY ./datax/datax-${DATAX_VERSION} /app/bin/datax
COPY ${CONFIG_FILE} /app/config.yaml
COPY ./rules /app/rules

# 复制 Nginx 配置和启动脚本
COPY docker/nginx.conf /etc/nginx/nginx.conf
COPY docker/start.sh /app/start.sh

# 设置工作目录
WORKDIR /app

# 设置可执行权限
RUN chmod +x /app/${BINARY_NAME} /app/bin/datax /app/start.sh

# 暴露端口
EXPOSE 80

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost/datax/api/v1/ping || exit 1

# 启动命令
CMD ["sh", "/app/start.sh"]