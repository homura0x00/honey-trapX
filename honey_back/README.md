# honey_back

蜜罐管理平台后端（Go 单体）。v1 重写，设计决策见 [docs/architecture.md](docs/architecture.md)。

## 技术栈

- Go + Gin（HTTP）
- MySQL + GORM（持久化；gorm AutoMigrate 建表，`internal/models` 为 schema 事实源，无 SQL 迁移）
- Redis（登录会话）
- Docker（Moby client，管理蜜罐容器实例）

## 目录

```
cmd/server/main.go   入口：装配 config → db → AutoMigrate → seed → feature → gin
internal/
  config/            配置加载（无全局变量）
  db/                MySQL/Redis 连接；Migrate/EnsureDatabase（AutoMigrate 与建库 flag）
  models/            手写 gorm 模型（schema 事实源）+ Principal 登录态
  auth/              登录/登出/会话 + Auth/Admin 中间件
  honey/             镜像模板 + 蜜罐部署闭环（seed.go：内置镜像模板）
  docker/            Docker Manager（唯一接触 daemon 的包）
  server/            gin 装配 + CORS
  utils/ utils/res/  bcrypt、统一响应/错误码
docs/architecture.md 架构决策记录
spec/                v1 讨论记录
```

## 运行

前置：MySQL、Redis、Docker 可用。部署分三步，前两步只需执行一次：

```bash
# 1) 创建数据库（连 MySQL 实例执行 CREATE DATABASE IF NOT EXISTS）
go run ./cmd/server -create-db

# 2) 创建默认管理员（users 表为空时创建；账户密码取 settings.yaml 的 admin 段，
#    默认 admin / admin123456，登录后请尽快修改）
go run ./cmd/server -create-admin

# 3) 启动服务（每次启动自动执行 gorm AutoMigrate 幂等建表；
#    镜像模板表为空时自动写入内置 seed）
go run ./cmd/server
```

## v1 接口（前缀 /api，均需登录 Cookie）

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | /auth/login | 登录 `{account,password}` |
| POST | /auth/logout | 登出 |
| GET | /auth/me | 当前登录用户 |
| GET | /images | 蜜罐镜像模板列表（seed 数据） |
| POST | /deployments | 创建部署 `{name,image_id}` |
| GET | /deployments[/:id] | 部署列表 / 详情 |
| POST | /deployments/:id/start | 启动 |
| POST | /deployments/:id/stop | 停止 |
| DELETE | /deployments/:id | 删除 |

响应统一 `{"code":0,"data":...,"message":""}`；code≠0 为业务错误（见
`internal/utils/res/errorcode.go`）。

## 明确不做（v1）

攻击捕获（attack_logs）、网络扫描、告警、攻击者画像、统计 dashboard、
system_logs、多节点/agent —— 均留到后续切片，见 docs/architecture.md。
