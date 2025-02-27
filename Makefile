# DataX Admin Makefile

# 变量定义
BINARY_NAME=datax-admin
MAIN_PACKAGE=./backend
DIST_DIR=./bin
FRONTEND_DIR=./frontend
DOCKER_IMAGE=datax-admin
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d %H:%M:%S')
COMMIT_HASH=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.CommitHash=$(COMMIT_HASH)'"

# Go 环境变量
export GO111MODULE=on
export CGO_ENABLED=0

# 颜色定义
BLUE=\033[0;34m
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

# 默认目标
.PHONY: all
all: clean build

# 帮助信息
.PHONY: help
help:
	@echo -e "$(BLUE)DataX Admin Makefile$(NC)"
	@echo ""
	@echo "使用方法:"
	@echo "  make [目标]"
	@echo ""
	@echo "目标:"
	@echo "  all             清理并构建项目 (默认)"
	@echo "  build           构建后端可执行文件 (当前平台)"
	@echo "  build-frontend  构建前端"
	@echo "  clean           清理构建产物"
	@echo "  test            运行测试"
	@echo "  docker          构建Docker镜像 (AMD64)"
	@echo "  docker-arm64    构建Docker镜像 (ARM64)"
	@echo "  docker-all      构建所有架构的Docker镜像"
	@echo "  linux           构建Linux可执行文件"
	@echo "  darwin          构建macOS可执行文件"
	@echo "  linux-amd64     构建Linux AMD64可执行文件"
	@echo "  linux-arm64     构建Linux ARM64可执行文件"
	@echo "  darwin-amd64    构建macOS AMD64可执行文件"
	@echo "  darwin-arm64    构建macOS ARM64可执行文件"
	@echo "  run             构建并运行应用"
	@echo "  help            显示帮助信息"
	@echo ""
	@echo "示例:"
	@echo "  make build"
	@echo "  make docker VERSION=v1.0.0"
	@echo "  make docker-arm64 VERSION=v1.0.0"
	@echo "  make docker-all VERSION=v1.0.0"
	@echo ""

# 创建目录
$(DIST_DIR):
	mkdir -p $(DIST_DIR)

# 清理
.PHONY: clean
clean:
	@echo -e "$(YELLOW)清理构建产物...$(NC)"
	rm -rf $(DIST_DIR)
	@echo -e "$(GREEN)清理完成$(NC)"

# 测试
.PHONY: test
test:
	@echo -e "$(YELLOW)运行测试...$(NC)"
	cd $(MAIN_PACKAGE) && go test -v ./...
	@echo -e "$(GREEN)测试完成$(NC)"

# 构建前端
.PHONY: build-frontend
build-frontend:
	@echo -e "$(YELLOW)构建前端...$(NC)"
	cd $(FRONTEND_DIR) && \
	if command -v pnpm > /dev/null; then \
		pnpm install && pnpm run build; \
	else \
		npm install && npm run build; \
	fi
	@echo -e "$(GREEN)前端构建完成$(NC)"

# 构建后端 (当前平台)
.PHONY: build
build: $(DIST_DIR)
	@echo -e "$(YELLOW)构建后端 (当前平台)...$(NC)"
	cd $(MAIN_PACKAGE) && go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME) .
	@echo -e "$(GREEN)后端构建完成: $(DIST_DIR)/$(BINARY_NAME)$(NC)"

# 构建 Linux AMD64
.PHONY: linux-amd64
linux-amd64: $(DIST_DIR)
	@echo -e "$(YELLOW)构建 Linux AMD64...$(NC)"
	cd $(MAIN_PACKAGE) && GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .
	@echo -e "$(GREEN)构建完成: $(DIST_DIR)/$(BINARY_NAME)-linux-amd64$(NC)"

# 构建 Linux ARM64
.PHONY: linux-arm64
linux-arm64: $(DIST_DIR)
	@echo -e "$(YELLOW)构建 Linux ARM64...$(NC)"
	cd $(MAIN_PACKAGE) && GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 .
	@echo -e "$(GREEN)构建完成: $(DIST_DIR)/$(BINARY_NAME)-linux-arm64$(NC)"

# 构建 macOS AMD64
.PHONY: darwin-amd64
darwin-amd64: $(DIST_DIR)
	@echo -e "$(YELLOW)构建 macOS AMD64...$(NC)"
	cd $(MAIN_PACKAGE) && GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@echo -e "$(GREEN)构建完成: $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64$(NC)"

# 构建 macOS ARM64
.PHONY: darwin-arm64
darwin-arm64: $(DIST_DIR)
	@echo -e "$(YELLOW)构建 macOS ARM64...$(NC)"
	cd $(MAIN_PACKAGE) && GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@echo -e "$(GREEN)构建完成: $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64$(NC)"

# 构建所有 Linux 平台
.PHONY: linux
linux: linux-amd64 linux-arm64
	@echo -e "$(GREEN)所有 Linux 平台构建完成$(NC)"

# 构建所有 macOS 平台
.PHONY: darwin
darwin: darwin-amd64 darwin-arm64
	@echo -e "$(GREEN)所有 macOS 平台构建完成$(NC)"

# 构建 Docker 镜像
.PHONY: docker
docker: linux-amd64
	@echo -e "$(YELLOW)构建 Docker 镜像 $(DOCKER_IMAGE):$(VERSION)...$(NC)"
	docker build --platform linux/amd64 -f Dockerfile \
		--build-arg BINARY_NAME=$(BINARY_NAME) \
		--build-arg BINARY_VERSION=linux-amd64 \
		--build-arg DATAX_VERSION=linux-amd64 \
		--build-arg CONFIG_FILE=./backend/config.yaml \
		-t $(DOCKER_IMAGE):$(VERSION)-amd64 .
	@echo -e "$(GREEN)Docker 镜像构建完成: $(DOCKER_IMAGE):$(VERSION)-amd64$(NC)"
	docker save $(DOCKER_IMAGE):$(VERSION)-amd64 > $(DOCKER_IMAGE).tar
	@echo -e "$(GREEN)Docker 镜像保存完成: $(DOCKER_IMAGE).tar$(NC)"

# 构建 Docker 镜像 (ARM64)
.PHONY: docker-arm64
docker-arm64: linux-arm64
	@echo -e "$(YELLOW)构建 ARM64 Docker 镜像 $(DOCKER_IMAGE):$(VERSION)-arm64$(NC)"
	docker build -f Dockerfile \
		--build-arg BINARY_NAME=$(BINARY_NAME) \
		--build-arg BINARY_VERSION=linux-arm64 \
		--build-arg DATAX_VERSION=linux-arm64 \
		--build-arg CONFIG_FILE=./backend/config.yaml \
		-t $(DOCKER_IMAGE):$(VERSION)-arm64 .
	@echo -e "$(GREEN)ARM64 Docker 镜像构建完成: $(DOCKER_IMAGE):$(VERSION)-arm64$(NC)"
	docker run -it -p 28080:80 $(DOCKER_IMAGE):$(VERSION)-arm64
	@echo -e "$(GREEN)Docker 容器已启动: $(DOCKER_IMAGE):$(VERSION)-arm64$(NC)"
	@echo "访问地址: http://localhost:28080/datax/"


# 运行应用
.PHONY: run
run: build
	@echo -e "$(YELLOW)运行应用...$(NC)"
	$(DIST_DIR)/$(BINARY_NAME)
