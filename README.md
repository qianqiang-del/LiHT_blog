# LiHT Blog

一个前后端分离的个人博客系统，采用 Go 后端 + Vue 3 前端架构，支持 Docker 一键部署。

## 📖 项目简介

LiHT Blog 是一个功能完整的博客平台，包含**博客前台**和**后台管理**两套前端应用。前台面向读者，提供文章浏览、分类/标签筛选、搜索、评论、点赞等功能；后台面向博主，提供文章管理、分类管理、标签管理、评论管理、用户管理、数据概览等能力。

## 🛠 技术栈

### 后端

| 技术 | 说明 |
|------|------|
| **Go 1.24** | 主语言 |
| **Gin** | HTTP Web 框架 |
| **GORM** | ORM 框架 |
| **MySQL 8.0** | 关系型数据库 |
| **Redis** | 缓存 / 计数（文章热度、点赞等） |
| **JWT** | 用户认证 |
| **Zap** | 高性能日志 |
| **Viper** | 配置管理 |
| **阿里云 OSS** | 图片/文件存储 |

### 前端

| 技术 | 说明 |
|------|------|
| **Vue 3** | 渐进式前端框架 |
| **TypeScript** | 类型安全 |
| **Vite 6** | 构建工具 |
| **Vue Router 4** | 路由管理 |
| **Marked** | Markdown 渲染 |

### 部署

| 技术 | 说明 |
|------|------|
| **Docker / Docker Compose** | 容器化部署 |
| **Nginx** | 静态资源托管 & API 反向代理 |

## 📁 项目结构

```
Blog/
├── cmd/server/              # 程序入口
├── internal/                # 业务代码（不对外暴露）
│   ├── api/                 # 路由 & 控制器
│   │   └── v1/              # API v1 版本
│   │       ├── article/     # 文章模块
│   │       ├── auth/        # 认证模块（登录/注册）
│   │       ├── author/      # 作者信息模块
│   │       ├── category/    # 分类模块
│   │       ├── comment/     # 评论模块
│   │       ├── dashboard/   # 后台概览模块
│   │       ├── tag/         # 标签模块
│   │       ├── upload/      # 文件上传模块
│   │       └── user/        # 用户管理模块
│   ├── app/                 # 应用初始化
│   ├── middleware/           # 中间件（JWT 认证、CORS、日志、Recovery）
│   ├── model/               # 数据模型
│   │   ├── dto/             # 请求/响应 DTO
│   │   └── entity/          # 数据库实体
│   ├── repository/          # 数据访问层（Repository 模式）
│   ├── service/             # 业务逻辑层
│   ├── stream/              # 消息队列消费（热度计算）
│   └── task/                # 定时任务
├── pkg/                     # 公共工具包
│   ├── config/              # 配置加载
│   ├── database/            # 数据库连接
│   ├── email/               # 邮件发送
│   ├── errors/              # 错误定义
│   ├── jwt/                 # JWT 工具
│   ├── logger/              # 日志初始化
│   ├── response/            # 统一响应封装
│   └── utils/               # 通用工具函数
├── frontend/                # 前端代码
│   ├── web/                 # 博客前台（端口 84）
│   └── admin/               # 后台管理（端口 85）
├── migrations/              # 数据库迁移脚本
├── config.yaml              # 本地配置文件（已 gitignore）
├── config.docker.yaml       # Docker 环境配置
├── docker-compose.yml       # Docker Compose 编排
├── nginx.conf               # Nginx 配置
├── Dockerfile               # 后端镜像构建
└── Makefile                 # 构建/运行命令
```

## ✨ 功能特性

### 博客前台

- 📝 文章列表 & 详情浏览（Markdown 渲染）
- 🔍 全文搜索
- 📂 按分类 / 标签筛选文章
- 💬 文章评论 & 回复
- ❤️ 文章 / 评论点赞
- 🔥 热门文章排行
- 👤 用户注册 / 登录
- 📄 关于页面（作者信息）

### 后台管理

- 📊 数据概览仪表盘（文章数、评论数、用户数等）
- 📝 文章发布 / 编辑 / 删除（支持封面上传、Markdown 编辑）
- 📂 分类管理（增删改查）
- 🏷 标签管理（增删改查）
- 💬 评论管理（ 删除）
- 👥 用户管理
- ⚙️ 个人设置（修改账号密码）
- 🔐 JWT 认证 & 权限控制

## 🚀 快速开始

### 环境要求

- Go 1.24+
- Node.js 18+
- MySQL 8.0
- Redis

### 1. 克隆项目

```bash
git clone <repo-url>
cd Blog
```

### 2. 配置

复制并编辑配置文件：

```bash
cp .env.example config.yaml
```

修改 `config.yaml` 中的数据库、Redis、JWT、邮件等配置项。

### 3. 启动后端

```bash
# 下载依赖
go mod tidy

# 运行
go run cmd/server/main.go

# 或使用 Makefile
make run
```

### 4. 启动前端

```bash
# 博客前台（端口 84）
cd frontend/web
npm install
npm run dev

# 后台管理（端口 85）
cd frontend/admin
npm install
npm run dev
```

### 5. 访问

| 服务 | 地址 |
|------|------|
| 博客前台 | http://47.122.109.1:84 |
| 后台管理 | http://47.122.109.1:85 |

## 🐳 Docker 部署

使用 Docker Compose 一键启动所有服务：

```bash
# 构建前端
cd frontend/web && npm run build && cd ../..
cd frontend/admin && npm run build && cd ../..

# 启动所有服务
docker-compose up -d
```

Docker Compose 会启动以下容器：

| 容器 | 端口 | 说明 |
|------|------|------|
| `blog-backend` | 9900 | Go API 服务 |
| `blog-mysql` | 3306 | MySQL 数据库 |
| `blog-redis` | 6379 | Redis 缓存 |
| `blog-nginx` | 84 / 85 | 前端静态资源 & 反向代理 |

## 📜 Makefile 命令

```bash
make build          # 编译项目
make run            # 运行项目
make test           # 运行测试
make test-coverage  # 生成测试覆盖率报告
make clean          # 清理构建产物
make deps           # 下载依赖
make fmt            # 格式化代码
make vet            # 静态检查
make lint           # 代码检查（需安装 golangci-lint）
make help           # 查看所有命令
```

## 📄 License

MIT
