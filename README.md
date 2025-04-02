# JWT权限系统

这是一个基于JWT的权限管理系统，提供用户认证、角色管理和权限控制功能。

## 功能特性

- 用户管理：注册、登录、信息管理
- 角色管理：创建角色、分配权限
- 权限管理：基于RBAC模型的权限控制
- JWT认证：生成令牌、验证令牌、刷新令牌
- 中间件：权限校验中间件

## 技术栈

- 后端：Go、Gin、GORM
- 数据库：PostgreSQL
- 认证：JWT

## 项目结构

```
jwt-auth-system/
├── cmd/                # 应用入口
│   └── main.go        # 主程序
├── config/            # 配置文件
│   └── config.yaml    # 配置文件
├── internal/          # 内部包
│   ├── config/        # 配置结构
│   ├── handler/       # HTTP处理器
│   ├── middleware/    # 中间件
│   ├── model/         # 数据模型
│   ├── repository/    # 数据访问层
│   └── service/       # 业务逻辑层
├── pkg/               # 公共包
│   └── auth/          # 认证工具
├── go.mod             # Go模块文件
└── README.md          # 项目说明
```

## 安装与运行

1. 克隆项目
2. 配置数据库
3. 运行应用

```bash
go run cmd/main.go
```

## API文档

### 认证API

- POST /api/auth/register - 用户注册
- POST /api/auth/login - 用户登录
- POST /api/auth/refresh - 刷新令牌
- GET /api/auth/profile - 获取用户信息

### 用户管理API

- GET /api/users - 获取用户列表
- GET /api/users/:id - 获取用户详情
- PUT /api/users/:id - 更新用户信息
- DELETE /api/users/:id - 删除用户

### 角色管理API

- GET /api/roles - 获取角色列表
- POST /api/roles - 创建角色
- PUT /api/roles/:id - 更新角色
- DELETE /api/roles/:id - 删除角色

### 权限管理API

- GET /api/permissions - 获取权限列表
- POST /api/roles/:id/permissions - 分配权限到角色