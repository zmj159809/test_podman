# 非root环境中使用Podman部署业务：问题与解决方案

## 环境信息
- **操作系统：** Fedora Remix for WSL2
- **用户权限：** 非root用户（rootless模式）
- **容器引擎：** Podman 5.5.2
- **编排工具：** podman-compose 1.4.1

---

## 主要问题及解决方案

### 1. Docker Compose工具缺失
**问题现象：**
```bash
bash: docker-compose: 未找到命令
```

**解决方案：**
```bash
sudo dnf install -y podman-compose
```

**说明：** 在Fedora系统中，需要使用 `podman-compose` 替代 `docker-compose`

---

### 2. 网络配置问题 - nftables错误
**问题现象：**
```bash
Error: unable to start container: netavark: nftables error: "nft" did not return successfully while applying ruleset
```

**根本原因：**
- WSL2环境中的nftables实现与Linux原生环境存在差异
- netavark网络后端在应用防火墙规则时失败
- 这是WSL2 + Podman + netavark的已知兼容性问题

**解决方案对比：**

| 方案 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| **host网络模式** | 性能最佳、配置简单 | 无网络隔离、端口冲突风险 | 开发环境、单服务 |
| **slirp4netns网络模式** | 网络隔离、端口映射、用户空间安全 | 性能略低 | 生产环境、多服务 |

**最终采用的解决方案：**
```yaml
# docker-compose.yml
version: '3.8'
services:
  web-server:
    build: .
    network_mode: slirp4netns  # 关键配置
    ports:
      - "15667:5667"
```

**依赖安装：**
```bash
sudo dnf install -y slirp4netns
```

---

### 3. PATH环境变量污染
**问题现象：**
系统PATH中包含大量Windows路径（`/mnt/c/*`），可能导致命令冲突

**解决方案：**
```bash
# 临时清理
export PATH=$(echo $PATH | tr ':' '\n' | grep -v '^/mnt/c' | tr '\n' ':' | sed 's/:$//')

# 永久生效 - 添加到 ~/.bashrc
echo 'PATH=$(echo $PATH | tr ":" "\n" | grep -v "^/mnt/c" | tr "\n" ":" | sed "s/:$//")'  >> ~/.bashrc
```

---

### 4. 容器镜像管理
**问题现象：**
构建过程中产生大量悬空镜像和临时容器

**解决方案：**
```bash
# 清理悬空镜像
podman image prune -f

# 系统全面清理
podman system prune -a -f

# 回收空间：393.1MB
```

---

## 网络模式详细对比

### netavark（默认）vs slirp4netns

| 特性 | netavark | slirp4netns |
|------|----------|-------------|
| **性能** | 高 | 中等 |
| **网络隔离** | 完整 | 完整 |
| **WSL2兼容性** | ❌ 存在问题 | ✅ 良好 |
| **端口映射** | 支持 | 支持 |
| **用户空间** | 部分 | 完全 |
| **配置复杂度** | 低 | 低 |

---

## 最终工作配置

**docker-compose.yml：**
```yaml
version: '3.8'
services:
  web-server:
    build: .
    network_mode: slirp4netns
    ports:
      - "15667:5667"
    environment:
      - PORT=5667
      - READ_TIMEOUT=15
      - WRITE_TIMEOUT=15
      - IDLE_TIMEOUT=60
      - LOG_LEVEL=info
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:5667/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 5s
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 128M
        reservations:
          cpus: '0.25'
          memory: 64M
```

**启动命令：**
```bash
podman-compose up -d
```

---

## 常用管理命令

```bash
# 查看服务状态
podman-compose ps

# 查看日志
podman-compose logs web-server

# 停止服务
podman-compose down

# 重启服务
podman-compose restart

# 实时查看日志
podman-compose logs -f web-server
```

---

## 关键经验教训

1. **WSL2环境特殊性：** 需要特别注意网络栈的兼容性问题
2. **网络后端选择：** slirp4netns是WSL2环境下的最佳选择
3. **工具链替换：** 使用 `podman-compose` 而不是 `docker-compose`
4. **环境隔离：** 及时清理PATH中的Windows路径污染
5. **资源管理：** 定期清理容器镜像避免空间浪费

### 最终效果
- ✅ **网络隔离：** 容器运行在独立的网络空间（10.0.2.0/24）
- ✅ **端口映射：** 外部15667端口映射到容器5667端口
- ✅ **健康检查：** 自动监控服务状态
- ✅ **资源限制：** CPU和内存资源受限
- ✅ **日志管理：** 通过 `podman-compose logs` 查看

## 故障排除指南

### 网络问题诊断
```bash
# 查看网络后端
podman info --format json | jq '.host.networkBackend'

# 检查网络列表
podman network ls

# 检查网络详情
podman network inspect <network_name>
```

### 容器状态检查
```bash
# 查看容器详细信息
podman inspect <container_name>

# 查看容器日志
podman logs <container_name>

# 进入容器调试
podman exec -it <container_name> /bin/sh
```

---

## 结论

这套解决方案既保证了容器的网络隔离，又避免了WSL2环境中的兼容性问题，是非root环境下部署容器化应用的可靠方案。通过使用slirp4netns网络模式，我们成功解决了WSL2环境下的网络配置问题，实现了稳定的容器化部署。
