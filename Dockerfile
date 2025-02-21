# 第一阶段：构建后端
FROM golang:1.22.2-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend .
RUN go build -o datax-admin .


# 第三阶段：最终镜像
FROM caddy:2.7-alpine
WORKDIR /usr/share/caddy


# 设置时区
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

# 复制前端构建产物
COPY frontend/dist ./dist/

# 复制后端二进制和配置
COPY --from=backend-builder /app/backend/datax-admin /app/datax-admin
COPY --from=backend-builder /app/backend/config.yaml /app/config.yaml
COPY --from=backend-builder /app/backend/bin/datax_linux_arm64 /app/bin/datax

# 复制 Caddy 配置
COPY Caddyfile /etc/caddy/Caddyfile
COPY start.sh /app/start.sh

# 创建启动脚本
WORKDIR /app
RUN caddy add-package github.com/caddyserver/transform-encoder && \
    caddy fmt --overwrite /etc/caddy/Caddyfile && \
    chmod +x start.sh

EXPOSE 80
CMD ["sh", "/app/start.sh"]
