# 包含 .env 文件中的变量
include .env
export

# 默认目标
.PHONY: help
help: ## 显示帮助信息
	@echo "可用的命令:"
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?## .*$$/ {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.PHONY: build
build: ## 构建镜像
	@echo "🔨 构建镜像 (版本: $(VERSION))..."
	podman compose build --no-cache

.PHONY: up
up: ## 启动服务
	@echo "▶️ 启动服务..."
	podman compose up -d

.PHONY: down
down: ## 停止服务
	@echo "🛑 停止服务..."
	podman compose down

.PHONY: restart
restart: down up ## 重启服务

.PHONY: deploy
deploy: down build up test ## 完整部署流程 (停止->构建->启动->测试)
	@echo "🎉 部署完成！版本: $(VERSION)"

.PHONY: redeploy
redeploy: deploy ## 重新部署 (deploy 的别名)

.PHONY: logs
logs: ## 查看日志
	@echo "📋 查看容器日志..."
	podman compose logs -f

.PHONY: status
status: ## 查看服务状态
	@echo "🔍 服务状态:"
	podman compose ps

.PHONY: test
test: ## 测试服务
	@echo "🧪 测试服务..."
	@if curl -f -s http://localhost:$(EXTERNAL_PORT)/health > /dev/null; then \
		echo "✅ 服务运行正常！"; \
		echo "🌐 访问地址: http://localhost:$(EXTERNAL_PORT)"; \
		echo "❤️ 健康检查: http://localhost:$(EXTERNAL_PORT)/health"; \
		echo "📊 指标: http://localhost:$(EXTERNAL_PORT)/metrics"; \
		echo "📋 版本信息:"; \
		curl -s http://localhost:$(EXTERNAL_PORT)/ | jq '.' 2>/dev/null || curl -s http://localhost:$(EXTERNAL_PORT)/; \
	else \
		echo "❌ 服务测试失败"; \
		make logs; \
		exit 1; \
	fi

.PHONY: clean
clean: down ## 清理容器和镜像
	@echo "🧹 清理容器和镜像..."
	podman container prune -f
	podman image prune -f

.PHONY: version
version: ## 显示当前版本
	@echo "📦 当前版本: $(VERSION)"
	@echo "🐹 Go 版本: $(GO_VERSION)"
	@echo "🔧 外部端口: $(EXTERNAL_PORT)"
	@echo "🔧 内部端口: $(INTERNAL_PORT)"

.PHONY: dev
dev: ## 开发模式 (构建并启动，显示日志)
	@echo "🛠️ 开发模式启动..."
	make deploy
	make logs

.PHONY: shell
shell: ## 进入容器 shell
	@echo "🐚 进入容器..."
	podman exec -it $$(podman compose ps -q web-server) /bin/sh

.PHONY: check-env
check-env: ## 检查环境变量
	@echo "🔍 检查环境变量:"
	@echo "VERSION=$(VERSION)"
	@echo "GO_VERSION=$(GO_VERSION)"
	@echo "EXTERNAL_PORT=$(EXTERNAL_PORT)"
	@echo "INTERNAL_PORT=$(INTERNAL_PORT)"
	@echo "READ_TIMEOUT=$(READ_TIMEOUT)"
	@echo "WRITE_TIMEOUT=$(WRITE_TIMEOUT)"
	@echo "IDLE_TIMEOUT=$(IDLE_TIMEOUT)"
	@echo "LOG_LEVEL=$(LOG_LEVEL)"

.PHONY: update-version
update-version: ## 更新版本号 (使用: make update-version VERSION=1.0.3)
	@if [ -z "$(NEW_VERSION)" ]; then \
		echo "❌ 请指定新版本号: make update-version NEW_VERSION=1.0.3"; \
		exit 1; \
	fi
	@echo "📝 更新版本号从 $(VERSION) 到 $(NEW_VERSION)..."
	@sed -i 's/VERSION=$(VERSION)/VERSION=$(NEW_VERSION)/' .env
	@echo "✅ 版本号已更新到 $(NEW_VERSION)"

.PHONY: quick-deploy
quick-deploy: ## 快速部署 (仅重启，不重新构建)
	@echo "⚡ 快速部署..."
	make restart
	make test

.PHONY: full-clean
full-clean: clean ## 完全清理 (包括停止所有相关容器)
	@echo "🧹 完全清理..."
	podman stop $$(podman ps -q --filter "ancestor=test_podman_web-server") 2>/dev/null || true
	podman rm $$(podman ps -aq --filter "ancestor=test_podman_web-server") 2>/dev/null || true
	podman rmi $$(podman images -q test_podman_web-server) 2>/dev/null || true

# 默认目标
.DEFAULT_GOAL := help
