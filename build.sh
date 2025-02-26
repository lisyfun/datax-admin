#!/bin/bash

# 设置错误时退出
set -e

# 设置默认值
GOOS=${GOOS:-linux}
GOARCH=${GOARCH:-amd64}
TAG=${TAG:-latest}
BUILD_FRONTEND=${BUILD_FRONTEND:-true}
BUILD_BACKEND=${BUILD_BACKEND:-true}
BUILD_DOCKER=${BUILD_DOCKER:-true}

echo "=== 开始构建 DataX-Admin ==="

# 显示构建信息
echo "开始构建 datax-admin:"
echo "操作系统: $GOOS"
echo "架构: $GOARCH"
echo "镜像标签: $TAG"

# 构建前端
if [ "$BUILD_FRONTEND" = "true" ]; then
    echo "构建前端..."
    cd frontend
    pnpm install
    pnpm run build
    cd ..
fi

# 构建后端
if [ "$BUILD_BACKEND" = "true" ]; then
    echo "构建后端..."
    cd backend

    # 创建bin目录（如果不存在）
    mkdir -p bin

    # 编译可执行文件
    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -o bin/datax-admin-$GOOS-$GOARCH

    # 创建符号链接
    ln -sf datax-admin-$GOOS-$GOARCH bin/datax-admin

    cd ..
fi

# 构建Docker镜像
if [ "$BUILD_DOCKER" = "true" ]; then
    echo "构建Docker镜像..."
    docker build --build-arg GOOS=$GOOS --build-arg GOARCH=$GOARCH -t datax-admin:$TAG .
fi

echo "构建完成！"
