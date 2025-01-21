#!/bin/bash

# 设置错误时退出
set -e

echo "=== 开始构建 DataX-Admin ==="

# 构建前端
echo "=== 开始构建前端 ==="
cd frontend

echo "构建前端项目..."
pnpm build

cd ..

# 构建后端
echo "=== 开始构建镜像 ==="
docker build --platform=linux/amd64 -f Dockerfile_linux_amd64 -t datax-admin:amd64 .
echo "=== 构建完成 ==="
