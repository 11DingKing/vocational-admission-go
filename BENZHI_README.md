# BENZHI_README

这是一个 Go 后端应用，提供年度计划、专业组、志愿投档、录取裁决、征集志愿和审计能力。

## 项目说明

- 项目：11DingKing/vocational-admission-go
- 项目用途：提供年度计划、专业组、志愿投档、录取裁决、征集志愿和审计能力。SQLite 迁移在启动时执行，所有请求经过会话鉴权、请求 ID 和统一错误处理。招生办使用 officer 角色，复核人员使用 reviewer 角色，admin 可执行全部操作。
- Go 工具链：`golang:1.23`
- 前端工具链：无

## 标准构建、运行和测试命令

进入容器后执行：

```bash
# 编译
cd '/app' && GOTOOLCHAIN=local go build ./...

# 启动
cd '/app' && GOTOOLCHAIN=local go run ./cmd/server

# 测试
cd '/app' && GOTOOLCHAIN=local go test ./...
```

## Docker 构建和进入容器

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh benzhi-task-66-amd64 linux/amd64
./build_benzhi_docker.sh benzhi-task-66-arm64 linux/arm64
docker run -it benzhi-task-66-amd64:latest
docker run -it --platform linux/arm64 benzhi-task-66-arm64:latest
```

## 题目验证命令

1. 预期退出码 1：`go test ./internal/service -run '^TestProcessBatchReportsEveryDecision$' -count=1`
