# GopherAI Makefile
# 使用方法: make <目标>
# 例如: make build    # 编译 Linux 二进制
#       make run      # 本地启动
#       make deploy   # 构建前后端，准备部署

APP_NAME = gopherai
BINARY = $(APP_NAME)
# Linux 编译参数（生产服务器目标）
LINUX_GOOS = linux
LINUX_GOARCH = amd64

.PHONY: build build-local run frontend-install frontend-build clean deploy help

help: ## 显示帮助信息
	@echo "GopherAI 构建工具"
	@echo ""
	@echo "用法: make <目标>"
	@echo ""
	@echo "目标:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## 编译 Linux 二进制（用于生产服务器）
	@echo "==> 编译 Linux amd64 二进制..."
	GOOS=$(LINUX_GOOS) GOARCH=$(LINUX_GOARCH) go build -ldflags="-s -w" -o $(BINARY) .
	@echo "==> 完成: $(BINARY)"
	@ls -lh $(BINARY)

build-local: ## 编译本地开发二进制
	@echo "==> 编译本地二进制..."
	go build -o $(BINARY) .
	@echo "==> 完成: $(BINARY)"

run: frontend-build build ## 本地启动完整项目
	@echo "==> 启动服务..."
	./$(BINARY)

frontend-install: ## 安装前端依赖
	@echo "==> 安装前端依赖..."
	cd vue-frontend && npm install

frontend-build: ## 构建前端生产包
	@echo "==> 构建前端..."
	cd vue-frontend && npm run build

clean: ## 清理构建产物
	@echo "==> 清理..."
	rm -f $(BINARY)
	rm -rf vue-frontend/dist/
	@echo "==> 清理完成"

deploy: build frontend-build ## 完整构建（后端 + 前端），准备部署
	@echo "==> 构建完成，可以部署了"
	@echo "后端二进制: $(BINARY)"
	@echo "前端静态文件: vue-frontend/dist/"
