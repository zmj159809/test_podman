# Test Podman Web服务器

这是一个经过全面优化的Go语言HTTP服务器，专为容器化部署而设计，支持统一的版本管理系统。

## 🚀 功能特性

### 核心功能
- **结构化设计**: 采用模块化架构，易于维护和扩展
- **配置管理**: 支持环境变量配置，具有合理的默认值
- **优雅关闭**: 支持SIGINT和SIGTERM信号的优雅关闭
- **请求日志**: 详细的请求日志记录，包括响应时间和状态码
- **安全头部**: 自动添加安全相关的HTTP头部

### API端点
- `GET /` - 主要端点，返回Hello World消息和服务器信息
- `GET /health` - 健康检查端点
- `GET /ready` - 就绪检查端点
- `GET /metrics` - 基础指标端点

### 容器优化
- **多阶段构建**: 使用Docker多阶段构建，最小化镜像大小
- **安全镜像**: 基于scratch镜像，使用非root用户运行
- **健康检查**: 内置Docker健康检查
- **资源限制**: 配置了CPU和内存限制

### 版本管理
- **统一版本控制**: 通过 `.env` 文件统一管理版本号
- **自动版本注入**: 构建时自动将版本信息编译进应用
- **镜像标签同步**: Docker镜像标签与应用版本保持一致
- **构建时间戳**: 记录应用的构建时间

## 📦 项目结构

```
test_podman/
├── main.go              # 主应用程序代码
├── go.mod              # Go模块定义
├── go.sum              # Go模块校验和
├── Dockerfile          # Docker构建文件
├── docker-compose.yml  # Docker Compose配置
├── .env               # 环境变量配置（包含版本号）
├── Makefile           # 自动化构建和部署命令
└── README.md          # 项目文档
```

## 🔧 配置选项

### 环境变量
| 变量名 | 默认值 | 描述 |
|--------|--------|------|
| PORT | 5667 | 服务器监听端口 |
| READ_TIMEOUT | 15 | 读取超时时间（秒） |
| WRITE_TIMEOUT | 15 | 写入超时时间（秒） |
| IDLE_TIMEOUT | 60 | 空闲超时时间（秒） |
| LOG_LEVEL | info | 日志级别 |

### 端口映射
- **内部端口**: 5667 (容器内)
- **外部端口**: 15667 (宿主机)

## 📚 版本管理

### 版本管理系统概述

本项目实现了统一的版本管理系统，通过修改 `.env` 文件中的 `VERSION` 变量，可以同时更新：
- Docker 镜像标签
- Go 应用程序内的版本信息
- 构建时间戳

### 版本管理命令

#### 1. 查看当前版本
```bash
make version
```

#### 2. 更新版本号
```bash
make update-version NEW_VERSION=1.0.4
```

#### 3. 部署新版本
```bash
make deploy
```

#### 4. 快速部署（仅重启，不重新构建）
```bash
make quick-deploy
```

### 完整工作流程

1. **修改代码**
2. **更新版本号**：
   ```bash
   make update-version NEW_VERSION=1.0.4
   ```
3. **部署新版本**：
   ```bash
   make deploy
   ```

### 可用命令

使用 `make help` 查看所有可用命令：

- `make build` - 构建镜像
- `make up` - 启动服务
- `make down` - 停止服务
- `make restart` - 重启服务
- `make deploy` - 完整部署流程（停止→构建→启动→测试）
- `make test` - 测试服务
- `make logs` - 查看日志
- `make status` - 查看服务状态
- `make clean` - 清理容器和镜像
- `make version` - 显示当前版本信息
- `make check-env` - 检查环境变量

### 版本信息访问

部署后，可以通过以下端点获取版本信息：

- **主页面**：`http://localhost:15667/`
- **健康检查**：`http://localhost:15667/health`
- **指标**：`http://localhost:15667/metrics`

所有端点都会返回包含版本号和构建时间的 JSON 响应。

### 工作原理

1. **版本注入**：在 Docker 构建过程中，通过 `ldflags` 将版本信息编译到 Go 二进制文件中
2. **环境变量**：Docker Compose 从 `.env` 文件读取版本信息
3. **镜像标签**：Docker 镜像使用版本号作为标签
4. **运行时显示**：Go 应用程序在运行时显示编译时注入的版本信息

### 注意事项

- 版本号格式建议使用语义化版本（如 1.0.0）
- 每次代码修改后都应该更新版本号
- 部署前请确保 `.env` 文件中的版本号是正确的
- 使用 `make deploy` 会自动进行完整的构建和测试流程

## 🛠️ 使用方法

### 本地开发
```bash
# 直接运行
go run main.go

# 编译后运行
go build -o test_podman
./test_podman
```

### Docker部署
```bash
# 构建镜像
docker build -t test_podman:latest .

# 运行容器
docker run -d --name test_podman -p 15667:5667 test_podman:latest
```

### Docker Compose部署
```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f web-server

# 停止服务
docker-compose down
```

## 🔍 API使用示例

### 主要端点
```bash
curl http://localhost:15667/
```
响应：
```json
{
  "message": "Hello World!",
  "timestamp": "2024-01-01T12:00:00Z",
  "version": "1.0.0",
  "uptime": "1h30m45s"
}
```

### 健康检查
```bash
curl http://localhost:15667/health
```
响应：
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "uptime": "1h30m45s",
  "version": "1.0.0"
}
```

### 就绪检查
```bash
curl http://localhost:15667/ready
```
响应：
```json
{
  "status": "ready",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### 指标端点
```bash
curl http://localhost:15667/metrics
```
响应：
```json
{
  "uptime_seconds": 5445.123,
  "timestamp": 1704110400,
  "version": "1.0.0",
  "go_version": "1.24"
}
```

## 🎯 性能优化

### 构建优化
- **静态编译**: 使用CGO_ENABLED=0进行静态编译
- **二进制优化**: 使用-ldflags='-w -s'减小二进制文件大小
- **依赖缓存**: Docker构建中优化依赖缓存层

### 运行时优化
- **HTTP超时**: 配置合理的读写超时时间
- **连接池**: 使用HTTP/1.1 keep-alive连接复用
- **日志优化**: 高效的结构化日志记录

### 容器优化
- **最小镜像**: 使用scratch作为基础镜像
- **非root用户**: 以普通用户权限运行
- **资源限制**: 配置CPU和内存限制
- **健康检查**: 内置健康检查机制

## 🔒 安全特性

### HTTP安全头部
- `X-Content-Type-Options: nosniff` - 防止MIME类型嗅探
- `X-Frame-Options: DENY` - 防止点击劫持
- `X-XSS-Protection: 1; mode=block` - XSS保护

### 容器安全
- **非root用户**: 使用专门的appuser用户运行
- **最小权限**: 仅包含必要的运行时文件
- **静态分析**: 无外部依赖的静态二进制文件

## 📈 监控和日志

### 日志格式
服务器输出详细的请求日志：
```
[SERVER] 2024/01/01 12:00:00 main.go:130: GET /health 200 1.234ms 127.0.0.1
```

### 监控指标
- **响应时间**: 每个请求的处理时间
- **状态码**: HTTP响应状态码
- **运行时间**: 服务器运行时间
- **请求来源**: 客户端IP地址

## 🛡️ 错误处理

### 优雅关闭
- 接收SIGINT和SIGTERM信号
- 30秒的优雅关闭超时
- 完成正在处理的请求后关闭

### 错误恢复
- 详细的错误日志记录
- 适当的HTTP状态码返回
- 防止应用程序崩溃

## 🔧 开发建议

### 代码规范
- 使用`go fmt`格式化代码
- 使用`go vet`检查代码
- 遵循Go语言最佳实践

### 测试
```bash
# 运行测试
go test -v ./...

# 性能测试
go test -bench=. ./...
```

### 部署检查
```bash
# 健康检查
curl -f http://localhost:15667/health

# 性能测试
ab -n 1000 -c 10 http://localhost:15667/
```

## 📝 更新日志

### v1.0.3 (当前版本)
- ✅ 统一版本管理系统
- ✅ 自动化 Makefile 部署命令
- ✅ 版本信息自动注入
- ✅ 镜像标签与应用版本同步
- ✅ 构建时间戳记录
- ✅ 一键版本更新和部署

### v1.0.0
- ✅ 完整的HTTP服务器实现
- ✅ 配置管理系统
- ✅ 优雅关闭机制
- ✅ 请求日志中间件
- ✅ 健康检查端点
- ✅ Docker容器优化
- ✅ 安全特性增强
- ✅ 中文注释和文档

---

**注意**: 此项目已经过全面优化，包括性能、安全性、可维护性和部署便利性。适合用于生产环境的容器化部署。
