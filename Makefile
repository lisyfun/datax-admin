# DataX Admin Makefile

# 变量定义
BINARY_NAME=datax-admin
MAIN_PACKAGE=./backend
DIST_DIR=./bin
FRONTEND_DIR=./frontend
DOCKER_IMAGE=datax-admin
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date '+%Y-%m-%d %H:%M:%S')
COMMIT_HASH=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.CommitHash=$(COMMIT_HASH)'"
# Docker Hub 用户名，可以通过环境变量设置或在命令行传入
DOCKER_USERNAME?=lisongyu

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
	@echo "  build           构建后端可执行文件 (Linux AMD64)"
	@echo "  build-arm64     构建后端可执行文件 (Linux ARM64)"
	@echo "  build-frontend  构建前端"
	@echo "  clean           清理构建产物"
	@echo "  run             运行应用（前端和后端）"
	@echo "  docker          构建Docker镜像 (AMD64)"
	@echo "  docker-arm64    构建Docker镜像 (ARM64)"
	@echo "  docker-all      构建所有架构的Docker镜像"
	@echo "  docker-push     推送多架构Docker镜像并创建latest标签"
	@echo "  docker-run      运行AMD64版本的Docker容器"
	@echo "  docker-run-arm64 运行ARM64版本的Docker容器"
	@echo "  help            显示帮助信息"
	@echo ""
	@echo "示例:"
	@echo "  make build"
	@echo "  make docker VERSION=v1.0.0"
	@echo "  make docker-arm64 VERSION=v1.0.0"
	@echo "  make docker-all VERSION=v1.0.0"
	@echo "  make docker-push DOCKER_USERNAME=lisongyu VERSION=v1.0.0"
	@echo "  make docker-run"
	@echo "  make docker-run-arm64"
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

# 运行应用（前端和后端）
.PHONY: run
run:
	@echo -e "$(YELLOW)启动前端服务...$(NC)"
	cd $(FRONTEND_DIR) && \
	if command -v pnpm > /dev/null; then \
		pnpm dev & \
	else \
		npm run dev & \
	fi
	@echo -e "$(YELLOW)启动后端服务...$(NC)"
	cd $(MAIN_PACKAGE) && go run main.go
	@echo -e "$(GREEN)服务已启动$(NC)"
	@echo -e "$(GREEN)前端访问地址: http://localhost:3000$(NC)"
	@echo -e "$(GREEN)后端访问地址: http://localhost:28080$(NC)"



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

# 构建后端 (AMD平台)
.PHONY: build
build: $(DIST_DIR)
	@echo -e "$(YELLOW)构建前端...$(NC)"
	cd $(FRONTEND_DIR) && \
	if command -v pnpm > /dev/null; then \
		pnpm install && pnpm run build; \
	else \
		npm install && npm run build; \
	fi
	@echo -e "$(GREEN)前端构建完成$(NC)"
	@echo -e "$(YELLOW)构建 Linux AMD64...$(NC)"
	cd $(MAIN_PACKAGE) && GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .
	@echo -e "$(GREEN)构建完成: $(DIST_DIR)/$(BINARY_NAME)-linux-amd64$(NC)"

# 构建 (当前平台)
.PHONY: build-arm64
build-arm64: $(DIST_DIR)
	@echo -e "$(YELLOW)构建前端...$(NC)"
	cd $(FRONTEND_DIR) && \
	if command -v pnpm > /dev/null; then \
		pnpm install && pnpm run build; \
	else \
		npm install && npm run build; \
	fi
	@echo -e "$(GREEN)前端构建完成$(NC)"
	@echo -e "$(YELLOW)构建 Linux AMD64...$(NC)"
	cd $(MAIN_PACKAGE) && GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 .
	@echo -e "$(GREEN)构建完成: $(DIST_DIR)/$(BINARY_NAME)-linux-arm64$(NC)"


# 构建 Docker 镜像
.PHONY: docker
docker: build
	@echo -e "$(YELLOW)构建 Docker 镜像 $(DOCKER_IMAGE):$(VERSION)...$(NC)"
	docker build --platform linux/amd64 -f Dockerfile \
		--build-arg BINARY_NAME=$(BINARY_NAME) \
		--build-arg BINARY_VERSION=linux-amd64 \
		--build-arg DATAX_VERSION=linux-amd64 \
		--build-arg CONFIG_FILE=./backend/config.yaml \
		-t $(DOCKER_IMAGE):$(VERSION)-amd64 .
	@echo -e "$(GREEN)Docker 镜像构建完成: $(DOCKER_IMAGE):$(VERSION)-amd64$(NC)"
	docker save $(DOCKER_IMAGE):$(VERSION)-amd64 > $(DOCKER_IMAGE)-$(VERSION)-amd64.tar
	@echo -e "$(GREEN)Docker 镜像保存完成: $(DOCKER_IMAGE)-$(VERSION)-amd64.tar$(NC)"

# 构建 Docker 镜像 (ARM64)
.PHONY: docker-arm64
docker-arm64: build-arm64
	@echo -e "$(YELLOW)构建 ARM64 Docker 镜像 $(DOCKER_IMAGE):$(VERSION)-arm64$(NC)"
	docker build -f Dockerfile \
		--build-arg BINARY_NAME=$(BINARY_NAME) \
		--build-arg BINARY_VERSION=linux-arm64 \
		--build-arg DATAX_VERSION=linux-arm64 \
		--build-arg CONFIG_FILE=./backend/config.yaml \
		-t $(DOCKER_IMAGE):$(VERSION)-arm64 .
	@echo -e "$(GREEN)ARM64 Docker 镜像构建完成: $(DOCKER_IMAGE):$(VERSION)-arm64$(NC)"
	# docker run -it -p 28080:80 $(DOCKER_IMAGE):$(VERSION)-arm64
	@echo -e "$(GREEN)Docker 容器已启动: $(DOCKER_IMAGE):$(VERSION)-arm64$(NC)"
	@echo "访问地址: http://localhost:28080/datax/"


# docker 启动 amd64
.PHONY: docker-run
docker-run:
	docker run -d \
		--name datax-admin \
		-p 28080:80 \
		-v $(pwd)/logs:/app/logs \
		$(DOCKER_IMAGE):$(VERSION)-amd64

# docker 启动 arm64
.PHONY: docker-run-arm64
docker-run-arm64:
	docker run -d \
		--name datax-admin \
		-p 28080:80 \
		-v $(pwd)/logs:/app/logs \
		$(DOCKER_IMAGE):$(VERSION)-arm64

# 构建所有架构的 Docker 镜像
.PHONY: docker-all
docker-all: docker docker-arm64
	@echo -e "$(GREEN)所有架构的 Docker 镜像构建完成$(NC)"

# 推送多架构 Docker 镜像并创建 latest 标签
.PHONY: docker-push
docker-push: docker-all
	@echo -e "$(YELLOW)推送多架构 Docker 镜像并创建 latest 标签...$(NC)"
	@echo -e "$(YELLOW)使用 Docker Hub 用户名: $(DOCKER_USERNAME)$(NC)"

	# 为镜像添加标签
	docker tag $(DOCKER_IMAGE):$(VERSION)-amd64 $(DOCKER_USERNAME)/$(DOCKER_IMAGE):$(VERSION)-amd64
	docker tag $(DOCKER_IMAGE):$(VERSION)-arm64 $(DOCKER_USERNAME)/$(DOCKER_IMAGE):$(VERSION)-arm64

	# 推送镜像到 Docker Hub
	docker push $(DOCKER_USERNAME)/$(DOCKER_IMAGE):$(VERSION)-amd64
	docker push $(DOCKER_USERNAME)/$(DOCKER_IMAGE):$(VERSION)-arm64

	# 删除现有的 latest 标签
	@echo -e "$(YELLOW)删除现有的 latest 标签...$(NC)"
	docker manifest rm $(DOCKER_USERNAME)/$(DOCKER_IMAGE):latest || true

	# 创建并推送多架构清单
	docker manifest create $(DOCKER_USERNAME)/$(DOCKER_IMAGE):latest --amend\
		$(DOCKER_USERNAME)/$(DOCKER_IMAGE):$(VERSION)-amd64 \
		$(DOCKER_USERNAME)/$(DOCKER_IMAGE):$(VERSION)-arm64

	# 为清单中的每个镜像添加架构信息
	docker manifest annotate $(DOCKER_USERNAME)/$(DOCKER_IMAGE):latest \
		$(DOCKER_USERNAME)/$(DOCKER_IMAGE):$(VERSION)-amd64 --os linux --arch amd64
	docker manifest annotate $(DOCKER_USERNAME)/$(DOCKER_IMAGE):latest \
		$(DOCKER_USERNAME)/$(DOCKER_IMAGE):$(VERSION)-arm64 --os linux --arch arm64

	# 推送清单
	docker manifest push $(DOCKER_USERNAME)/$(DOCKER_IMAGE):latest

	@echo -e "$(GREEN)多架构 Docker 镜像推送完成$(NC)"
	@echo -e "$(GREEN)镜像可通过以下方式拉取:$(NC)"
	@echo "  docker pull $(DOCKER_USERNAME)/$(DOCKER_IMAGE):latest"
