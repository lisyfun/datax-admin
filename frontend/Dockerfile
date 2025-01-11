# 构建阶段
FROM node:18-alpine as builder

# 设置工作目录
WORKDIR /app

# 复制项目文件
COPY package.json pnpm-lock.yaml ./

# 安装 pnpm 和依赖
RUN npm install -g pnpm && pnpm install

# 复制源代码
COPY . .

# 构建项目
RUN pnpm vite build

# 运行阶段
FROM caddy:2-alpine

# 复制 Caddy 配置
COPY Caddyfile /etc/caddy/Caddyfile

# 从构建阶段复制构建产物
COPY --from=builder /app/dist /usr/share/caddy

# 暴露端口
EXPOSE 80
EXPOSE 443

# Caddy 会自动启动并加载配置文件
