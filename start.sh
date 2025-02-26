#!/bin/sh

# 启动后端服务
/app/datax-admin &

# 等待后端服务启动
echo "等待后端服务启动..."
while ! curl -s http://localhost:28080/datax/api/ping > /dev/null; do
    sleep 1
done
echo "后端服务已启动"

# 启动 Caddy
echo "启动 Caddy 服务..."
caddy run --config /etc/caddy/Caddyfile --adapter caddyfile
