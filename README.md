# 职业本科多省招生协同平台

提供年度计划、专业组、志愿投档、录取裁决、征集志愿和审计能力。SQLite 迁移在启动时执行，所有请求经过会话鉴权、请求 ID 和统一错误处理。招生办使用 officer 角色，复核人员使用 reviewer 角色，admin 可执行全部操作。

```bash
go run ./cmd/server
curl http://localhost:8080/healthz
```
