# go-boot-swagger

[![Go Version](https://img.shields.io/github/go-mod/go-version/xudefa/go-boot-swagger)](https://go.dev/) [![License](https://img.shields.io/github/license/xudefa/go-boot-swagger)](./LICENSE) [![Build Status](https://img.shields.io/github/actions/workflow/status/xudefa/go-boot-swagger/test.yml?branch=master)](https://github.com/xudefa/go-boot-swagger/actions) [![Go Reference](https://pkg.go.dev/badge/github.com/xudefa/go-boot-swagger.svg)](https://pkg.go.dev/github.com/xudefa/go-boot-swagger) [![Go Report Card](https://goreportcard.com/badge/github.com/xudefa/go-boot-swagger)](https://goreportcard.com/report/github.com/xudefa/go-boot-swagger)

基于 [go-boot](https://github.com/xudefa/go-boot) 的 Swagger API 文档集成模块。提供 Swagger UI 配置、路由注册和安全定义功能，兼容 Gin 和 Hertz 框架，支持多适配器。

> 设计理念：遵循 go-boot 的开发规范，通过函数式选项模式配置 Swagger，支持多框架适配和自动配置。

## 整体架构

```
┌───────────────────────────────────────────────────────────────────────┐
│                    go-boot ApplicationContext                         │
│  ┌───────────┐ ┌──────────────┐ ┌───────────┐ ┌───────────┐           │
│  │ Container │ │  Environment │ │ Lifecycle │ │ EventBus  │           │
│  └───────────┘ └──────────────┘ └───────────┘ └───────────┘           │
│                       ┌─────────────────────┐                         │
│                       │ AutoConfig Registry │                         │
│                       └─────────────────────┘                         │
└───────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
                    ┌───────────────────────────────┐
                    │   go-boot-swagger Starter     │
                    │  ┌─────────────────────────┐  │
                    │  │ SwaggerConfig Bean      │  │
                    │  │ GinAdapter              │  │
                    │  │ HertzAdapter            │  │
                    │  │ UI Routes               │  │
                    │  │ Security Definitions    │  │
                    │  └─────────────────────────┘  │
                    └───────────────────────────────┘
```

## 目录

- [快速开始](#快速开始)
- [功能特性](#功能特性)
- [配置选项](#配置选项)
- [路由注册](#路由注册)
- [安全定义](#安全定义)
- [多框架适配](#多框架适配)
- [项目结构](#项目结构)
- [开发指南](#开发指南)
- [贡献](#贡献)
- [许可证](#许可证)

## 快速开始

### 安装

```bash
# 安装核心框架
go get github.com/xudefa/go-boot

# 安装 Swagger 集成模块
go get github.com/xudefa/go-boot-swagger
```

### 最小示例

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/xudefa/go-boot/boot"
    "github.com/xudefa/go-boot/core"
    "github.com/xudefa/go-boot-swagger/swagger"
)

func main() {
    app, err := boot.NewApplication(
        boot.WithAppName("my-api-app"),
        boot.WithVersion("1.0.0"),
    )
    if err != nil {
        panic(err)
    }
    defer app.Stop()

    // 创建 Gin Engine
    engine := gin.Default()

    // 注册 Swagger 路由
    swagger.Register(engine,
        swagger.WithTitle("My API"),
        swagger.WithDescription("API Documentation"),
        swagger.WithVersion("1.0"),
        swagger.WithHost("localhost:8080"),
    )

    // 注册到容器
    app.Container().Register("ginEngine", core.Bean(engine))

    // 启动应用
    app.Start()

    // 访问 Swagger UI: http://localhost:8080/swagger/index.html

    // 等待终止信号
    app.WaitForSignal()
}
```

## 功能特性

| 特性 | 说明 |
|------|------|
| Swagger UI | 集成 Swagger UI 界面和路由 |
| 函数式选项 | 灵活的配置（标题、描述、版本、主机等） |
| 多框架适配 | 支持 Gin 和 Hertz 框架的适配器 |
| 安全定义 | 支持 API Key、OAuth2 等安全定义 |
| CORS 支持 | 为 Swagger UI 提供 CORS 支持 |
| 健康检查 | 内置健康检查端点 |
| 自动配置 | 通过自动配置注册 SwaggerConfig Bean |
| 向后兼容 | 保留 Gin 框架的向后兼容 API |

## 配置选项

### 函数式选项

```go
swagger.Register(engine,
    swagger.WithEnabled(true),
    swagger.WithPath("/swagger/*any"),
    swagger.WithTitle("My API"),
    swagger.WithDescription("API Documentation"),
    swagger.WithVersion("1.0"),
    swagger.WithHost("localhost:8080"),
    swagger.WithBasePath("/api"),
    swagger.WithSchemes("http", "https"),
    swagger.WithContact("Developer", "dev@example.com", "https://example.com"),
    swagger.WithLicense("MIT", "https://opensource.org/licenses/MIT"),
)
```

### 安全定义

```go
// API Key 认证
swagger.Register(engine,
    swagger.WithSecurityDefinition("ApiKeyAuth", swagger.SecurityDefinition{
        Type:        "apiKey",
        Description: "API Key 认证",
        Name:        "X-API-Key",
        In:          "header",
    }),
)

// OAuth2 认证
swagger.Register(engine,
    swagger.WithSecurityDefinition("OAuth2", swagger.SecurityDefinition{
        Type:             "oauth2",
        Description:      "OAuth2 认证",
        AuthorizationURL: "https://example.com/oauth/authorize",
        TokenURL:         "https://example.com/oauth/token",
        Scopes: map[string]string{
            "read":  "读取权限",
            "write": "写入权限",
        },
    }),
)
```

## 路由注册

### Gin 框架

```go
import "github.com/xudefa/go-boot-swagger/swagger"

// 在 Engine 上注册
swagger.Register(engine,
    swagger.WithTitle("My API"),
)

// 在 RouterGroup 上注册
swagger.RegisterWithGroup(group,
    swagger.WithTitle("My API"),
)
```

### Hertz 框架

```go
import "github.com/xudefa/go-boot-swagger/swagger"

// 使用 Hertz 适配器
adapter := swagger.NewHertzAdapter()
adapter.Register(hertzEngine,
    swagger.WithTitle("My API"),
)
```

## 多框架适配

### 适配器模式

go-boot-swagger 采用适配器模式支持多框架：

```go
// Gin 适配器
ginAdapter := swagger.NewGinAdapter()
middleware := ginAdapter.Middleware(opts...)

// Hertz 适配器
hertzAdapter := swagger.NewHertzAdapter()
hertzMiddleware := hertzAdapter.Middleware(opts...)
```

### 中间件

```go
// Gin 中间件
engine.Use(swagger.Middleware(
    swagger.WithTitle("My API"),
))

// CORS 中间件
engine.Use(swagger.CORSForSwagger())
```

## 项目结构

```
go-boot-swagger/
├── swagger.go              # Swagger 配置和核心功能
├── adapter.go              # 适配器接口定义
├── adapter_gin.go          # Gin 框架适配器
├── adapter_hertz.go        # Hertz 框架适配器
├── register.go             # 路由注册辅助
├── autoconfig.go           # 自动配置注册
├── README.md
├── LICENSE
└── go.mod
```

## 开发指南

### 构建

```bash
go build ./...
```

### 测试

```bash
go test ./...
go test -cover ./...       # 带覆盖率
go test -race ./...        # 数据竞争检测
```

### 代码规范

```bash
go fmt ./...
golangci-lint run
```

## 贡献

欢迎提交 Issue 和 Pull Request！详细贡献指南请参阅 [CONTRIBUTING.md](./CONTRIBUTING.md)。

## 许可证

本项目采用 MIT 许可证 — 详情请参阅 [LICENSE](./LICENSE) 文件。