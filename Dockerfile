# 多阶段构建以获得最佳镜像大小
ARG GO_VERSION=1.24
FROM golang:${GO_VERSION}-alpine AS builder

# 接收构建参数
ARG VERSION=1.0.0
ARG CGO_ENABLED=0

# 创建非root用户（跳过包更新以避免网络问题）
RUN adduser -D -g '' appuser

# 设置工作目录
WORKDIR /build

# 先复制依赖文件以便更好地缓存
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download && go mod verify

# 复制源代码
COPY . .

# 使用优化选项构建应用程序，并将版本信息编译进二进制文件
RUN CGO_ENABLED=${CGO_ENABLED} GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static' -X 'test_podman/version.Version=${VERSION}' -X 'test_podman/version.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)'" \
    -a -installsuffix cgo \
    -o main .

# 最终阶段 - 最小运行时镜像
FROM scratch

# 复制用户信息
COPY --from=builder /etc/passwd /etc/passwd

# 复制二进制文件
COPY --from=builder /build/main /app/main

# 使用非root用户
USER appuser

# 暴露端口
EXPOSE 15667

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/main", "health"] || exit 1

# 设置环境变量
ENV PORT=5667
ENV READ_TIMEOUT=15
ENV WRITE_TIMEOUT=15
ENV IDLE_TIMEOUT=60
ENV LOG_LEVEL=info

# 运行应用程序
ENTRYPOINT ["/app/main"]

