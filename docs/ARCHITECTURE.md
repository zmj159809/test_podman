# 项目架构

本项目已重构为模块化结构，便于阅读、维护和扩展。

## 目录结构

```
├── config/          # 配置管理
│   └── config.go    # 环境变量配置加载
├── handlers/        # HTTP处理器
│   └── handlers.go  # 所有HTTP端点处理器
├── middleware/      # 中间件
│   └── middleware.go # 日志记录等中间件
├── server/          # 服务器管理
│   └── server.go    # HTTP服务器设置和生命周期
├── types/           # 类型定义
│   └── types.go     # 共享的结构体定义
├── version/         # 版本管理
│   └── version.go   # 版本信息管理
└── main.go          # 应用程序入口点
```

## 模块说明

### config/ - 配置管理
- **config.go**: 从环境变量加载应用程序配置
- 提供默认值和环境变量覆盖机制

### handlers/ - HTTP处理器
- **handlers.go**: 包含所有HTTP端点的处理逻辑
- 包括健康检查、指标、主页等端点
- 采用结构体方法模式，便于依赖注入

### middleware/ - 中间件
- **middleware.go**: HTTP中间件实现
- 日志记录中间件，记录请求详情
- 自定义ResponseWriter用于捕获状态码

### server/ - 服务器管理
- **server.go**: HTTP服务器的设置和生命周期管理
- 路由配置
- 优雅启动和关闭

### types/ - 类型定义
- **types.go**: 共享的数据结构定义
- 配置结构体、响应结构体等

### version/ - 版本管理
- **version.go**: 版本信息管理
- 支持构建时通过ldflags注入版本信息

### main.go - 应用程序入口
- 应用程序启动逻辑
- 信号处理和优雅关闭
- 模块组装和依赖注入

## 优势

1. **可维护性**: 每个模块职责明确，便于维护
2. **可测试性**: 模块化设计便于单元测试
3. **可扩展性**: 新功能可以轻松添加到相应模块
4. **可读性**: 代码结构清晰，便于理解
5. **复用性**: 模块可以在其他项目中复用

## 构建和部署

构建和部署流程保持不变，所有现有的Makefile命令仍然有效：

```bash
make build    # 构建Docker镜像
make deploy   # 完整部署
make test     # 测试服务
```
