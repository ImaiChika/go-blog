# Go Blog 项目说明

这是一个用于学习 Go Web 开发的一步步手敲博客项目。当前项目采用前后端分离结构：后端使用 Go + Gin + GORM，提供用户认证、文章管理、图片上传和 Redis 缓存；前端使用 Vue 3 + Vite，目前已经搭好基础路由和页面框架。

## 功能介绍

### 后端功能

- 用户注册：支持用户名、密码校验，密码使用 bcrypt 加密后入库。
- 用户登录：登录成功后返回 JWT Token，有效期为 24 小时。
- JWT 鉴权：创建、更新、删除文章和上传图片都需要在请求头中携带 `Authorization: Bearer <token>`。
- 文章列表：支持分页查询，默认 `page=1`、`page_size=10`，按创建时间倒序返回。
- 文章详情：先查 Redis 缓存，缓存未命中再查 MySQL，并把结果缓存 1 小时。
- 文章创建：登录用户可以创建文章，作者字段来自 Token 中的用户名。
- 文章更新：只有文章作者本人可以更新文章，更新后会删除对应 Redis 缓存。
- 文章删除：只有文章作者本人可以删除文章，删除后会删除对应 Redis 缓存。
- 图片上传：登录后可上传 jpg/png 图片，大小限制 5MB，文件保存到 `go-blog-backend/uploads`，并通过 `/uploads/<filename>` 静态路径访问。
- 参数校验：注册、登录、文章、分页、封面图 URL 等输入都有基础校验和中文错误提示。

### 前端功能

- Vue 3 + Vite 项目结构已经初始化。
- 已配置 Vue Router。
- 当前页面路由：
  - `/`：首页占位页。
  - `/login`：登录页占位页。
  - `/create`：发布文章页占位页。
  - `/posts/:id`：文章详情页占位页。
- 当前前端还没有接入真实后端接口，主要是基础页面骨架。

## 目录结构

```text
.
├── go-blog-backend/          # Go 后端服务
│   ├── config/               # MySQL、Redis、环境变量配置
│   ├── controllers/          # 注册登录、文章、上传接口逻辑
│   ├── middleware/           # JWT 鉴权中间件
│   ├── models/               # GORM 数据模型
│   ├── routes/               # Gin 路由定义
│   ├── uploads/              # 上传文件保存目录
│   ├── utils/                # JWT 工具函数
│   ├── main.go               # 后端入口
│   ├── go.mod
│   └── go.sum
├── go-blog-frontend/         # Vue 前端项目
│   ├── src/
│   │   ├── router/           # 前端路由
│   │   ├── views/            # 页面组件
│   │   ├── App.vue
│   │   └── main.js
│   ├── package.json
│   └── vite.config.js
└── test.http                 # REST Client 接口调试示例
```

## 技术栈

### 后端

- Go
- Gin
- GORM
- MySQL
- Redis
- JWT
- bcrypt

### 前端

- Vue 3
- Vue Router
- Vite

## 环境变量

后端启动时会读取当前工作目录下的 `.env` 文件；如果环境变量已经在系统里存在，则系统环境变量优先，不会被 `.env` 覆盖。

建议在 `go-blog-backend/.env` 中配置：

```env
GO_BLOG_DB_DSN=root:123456@tcp(127.0.0.1:3306)/go_blog?charset=utf8mb4&parseTime=True&loc=Local
GO_BLOG_REDIS_ADDR=localhost:6379
GO_BLOG_REDIS_PASSWORD=
GO_BLOG_REDIS_DB=0
GO_BLOG_JWT_SECRET=your-secret-key
```

如果不配置，项目会使用代码中的默认值。

## 后端启动方法

### 1. 准备 MySQL

先确保本机 MySQL 已启动，并创建数据库：

```sql
CREATE DATABASE go_blog DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

项目启动时会通过 GORM `AutoMigrate` 自动创建或更新 `posts`、`users` 表。

### 2. 准备 Redis

确保本机 Redis 已启动，默认连接地址是：

```text
localhost:6379
```

文章详情接口会使用 Redis 缓存，因此 Redis 不可用时后端会启动失败。

### 3. 启动后端

```bash
cd go-blog-backend
go mod tidy
go run .
```

启动成功后，后端服务监听：

```text
http://localhost:8080
```

## 前端启动方法

```bash
cd go-blog-frontend
npm install
npm run dev
```

Vite 会在终端输出本地访问地址，通常是：

```text
http://localhost:5173
```

当前前端页面仍是占位内容，后续可以继续把登录、文章列表、详情、发布和图片上传接入后端 API。

## 接口说明

### 注册

```http
POST /register
Content-Type: application/json

{
  "username": "imai",
  "password": "my_secure_password"
}
```

### 登录

```http
POST /login
Content-Type: application/json

{
  "username": "imai",
  "password": "my_secure_password"
}
```

返回：

```json
{
  "token": "<jwt-token>"
}
```

### 获取文章列表

```http
GET /api/v1/posts?page=1&page_size=10
```

返回字段包含：

- `data`：文章列表。
- `total`：文章总数。
- `page`：当前页码。
- `page_size`：每页数量。

### 获取文章详情

```http
GET /api/v1/posts/1
```

返回字段包含：

- `data`：文章详情。
- `source`：数据来源，可能是 `redis` 或 `db`。

### 创建文章

```http
POST /api/v1/posts
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "title": "新的文章",
  "content": "这是内容",
  "cover_image": "http://localhost:8080/uploads/example.png"
}
```

### 更新文章

```http
PUT /api/v1/posts/1
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "title": "更新后的标题",
  "content": "更新后的内容",
  "cover_image": "http://localhost:8080/uploads/example.png"
}
```

只有文章作者本人可以更新。

### 删除文章

```http
DELETE /api/v1/posts/1
Authorization: Bearer <jwt-token>
```

只有文章作者本人可以删除。

### 上传图片

```http
POST /api/v1/upload
Authorization: Bearer <jwt-token>
Content-Type: multipart/form-data
```

表单字段名必须是：

```text
file
```

限制：

- 只允许 jpg/png。
- 文件大小不能超过 5MB。

成功返回：

```json
{
  "message": "文件上传成功",
  "url": "http://localhost:8080/uploads/<filename>.png"
}
```

## 使用 test.http 调试接口

根目录下的 `test.http` 已经整理了一批接口请求示例，适合在 VS Code 的 REST Client 插件里直接运行。

推荐调试顺序：

1. 注册用户。
2. 登录获取 Token。
3. 用 Token 上传图片。
4. 用 Token 创建文章。
5. 查询文章列表和文章详情。
6. 测试更新、删除和无权限访问。

## 运行测试

后端包含配置、参数校验、上传和文章权限相关测试。

```bash
cd go-blog-backend
go test ./...
```

注意：部分文章控制器测试会初始化 MySQL 和 Redis，因此运行完整测试前需要保证数据库和 Redis 可用。

前端可以运行构建检查：

```bash
cd go-blog-frontend
npm run build
```

## 当前项目状态

- 后端核心接口已经具备博客项目的基础能力。
- 用户认证、文章作者权限、分页、上传文件校验、Redis 缓存失效都已有实现。
- 前端目前还是页面骨架，还没有把表单和列表接入后端。
- 数据库迁移使用的是 GORM `AutoMigrate`，还没有独立 migration 文件。
- 上传文件保存在本地目录，暂未接入对象存储或 CDN。

## 后续可以继续做的事情

- 前端接入登录接口，并保存 JWT。
- 首页接入文章列表和分页。
- 文章详情页渲染后端文章内容。
- 发布文章页接入图片上传和创建文章接口。
- 增加编辑文章页面。
- 增加退出登录、登录态判断和路由守卫。
- 为用户注册增加重复用户名的友好提示。
- 给后端补充更完整的错误处理和统一响应结构。
- 增加 `.gitignore`，避免 `.DS_Store`、本地 `.env`、`node_modules`、上传测试文件进入版本管理。
