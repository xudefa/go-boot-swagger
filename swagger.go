// Package swagger 提供 Swagger API 文档集成功能。
//
// 支持 Swagger UI 的配置和路由注册，兼容 Gin 和 Hertz 框架。
// 提供安全定义、CORS 支持和健康检查端点。
//
// 核心组件：
//   - Config: Swagger 配置，包含标题、描述、版本等
//   - SecurityDefinition: API 安全定义
//   - Option: 函数式配置选项
package swagger

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Config Swagger配置
type Config struct {
	Enabled             bool                          // 是否启用Swagger
	Path                string                        // Swagger路径
	Title               string                        // API标题
	Description         string                        // API描述
	Version             string                        // API版本
	Host                string                        // API主机地址
	BasePath            string                        // API基础路径
	Schemes             []string                      // 支持的协议
	ContactName         string                        // 联系人姓名
	ContactEmail        string                        // 联系人邮箱
	ContactURL          string                        // 联系人URL
	LicenseName         string                        // 许可证名称
	LicenseURL          string                        // 许可证URL
	SecurityDefinitions map[string]SecurityDefinition // 安全定义
}

// SecurityDefinition 安全定义
type SecurityDefinition struct {
	Type             string            // 安全类型：apiKey、oauth2等
	Description      string            // 描述
	Name             string            // 参数名称
	In               string            // 参数位置：header、query等
	AuthorizationURL string            // 授权URL
	TokenURL         string            // 令牌URL
	Scopes           map[string]string // 权限范围
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Enabled:     true,
		Path:        "/swagger/*any",
		Title:       "API Documentation",
		Description: "This is a sample server API.",
		Version:     "1.0",
		Host:        "localhost:8080",
		BasePath:    "/",
		Schemes:     []string{"http", "https"},
	}
}

// Option 配置选项函数类型
type Option func(*Config)

// WithEnabled 启用或禁用Swagger
func WithEnabled(enabled bool) Option {
	return func(c *Config) {
		c.Enabled = enabled
	}
}

// WithPath 设置Swagger路径
func WithPath(path string) Option {
	return func(c *Config) {
		c.Path = path
	}
}

// WithTitle 设置API标题
func WithTitle(title string) Option {
	return func(c *Config) {
		c.Title = title
	}
}

// WithDescription 设置API描述
func WithDescription(description string) Option {
	return func(c *Config) {
		c.Description = description
	}
}

// WithVersion 设置API版本
func WithVersion(version string) Option {
	return func(c *Config) {
		c.Version = version
	}
}

// WithHost 设置API主机
func WithHost(host string) Option {
	return func(c *Config) {
		c.Host = host
	}
}

// WithBasePath 设置API基础路径
func WithBasePath(basePath string) Option {
	return func(c *Config) {
		c.BasePath = basePath
	}
}

// WithSchemes 设置支持的协议
func WithSchemes(schemes ...string) Option {
	return func(c *Config) {
		c.Schemes = schemes
	}
}

// WithContact 设置联系信息
func WithContact(name, email, url string) Option {
	return func(c *Config) {
		c.ContactName = name
		c.ContactEmail = email
		c.ContactURL = url
	}
}

// WithLicense 设置许可证信息
func WithLicense(name, url string) Option {
	return func(c *Config) {
		c.LicenseName = name
		c.LicenseURL = url
	}
}

// WithSecurityDefinition 添加安全定义
func WithSecurityDefinition(name string, def SecurityDefinition) Option {
	return func(c *Config) {
		if c.SecurityDefinitions == nil {
			c.SecurityDefinitions = make(map[string]SecurityDefinition)
		}
		c.SecurityDefinitions[name] = def
	}
}

// WithSecurityDefinitions 批量设置安全定义
func WithSecurityDefinitions(defs map[string]SecurityDefinition) Option {
	return func(c *Config) {
		c.SecurityDefinitions = defs
	}
}

// UpdateConfig 从另一个配置更新当前配置
func (c *Config) UpdateConfig(other *Config) {
	if other != nil {
		c.Enabled = other.Enabled
		if other.Path != "" {
			c.Path = other.Path
		}
		if other.Title != "" {
			c.Title = other.Title
		}
		if other.Description != "" {
			c.Description = other.Description
		}
		if other.Version != "" {
			c.Version = other.Version
		}
		if other.Host != "" {
			c.Host = other.Host
		}
		if other.BasePath != "" {
			c.BasePath = other.BasePath
		}
		if other.Schemes != nil {
			c.Schemes = other.Schemes
		}
		if other.ContactName != "" {
			c.ContactName = other.ContactName
		}
		if other.ContactEmail != "" {
			c.ContactEmail = other.ContactEmail
		}
		if other.ContactURL != "" {
			c.ContactURL = other.ContactURL
		}
		if other.LicenseName != "" {
			c.LicenseName = other.LicenseName
		}
		if other.LicenseURL != "" {
			c.LicenseURL = other.LicenseURL
		}
		if other.SecurityDefinitions != nil {
			c.SecurityDefinitions = other.SecurityDefinitions
		}
	}
}

// Middleware 返回Swagger中间件（Gin框架，向后兼容）
func Middleware(opts ...Option) gin.HandlerFunc {
	adapter := NewGinAdapter()
	middleware := adapter.Middleware(opts...)
	if fn, ok := middleware.(gin.HandlerFunc); ok {
		return fn
	}
	return func(c *gin.Context) {
		c.Next()
	}
}

// Register 注册Swagger路由（Gin框架，向后兼容）
func Register(router *gin.Engine, opts ...Option) {
	RegisterForGin(router, opts...)
}

// RegisterWithGroup 在路由组中注册Swagger（Gin框架，向后兼容）
func RegisterWithGroup(group *gin.RouterGroup, opts ...Option) {
	RegisterWithGroupForGin(group, opts...)
}

// HealthCheck 健康检查端点（Gin框架，向后兼容）
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// CORSForSwagger 为Swagger UI提供CORS支持（Gin框架，向后兼容）
func CORSForSwagger() gin.HandlerFunc {
	adapter := NewGinAdapter()
	cors := adapter.CORSForSwagger()
	if fn, ok := cors.(gin.HandlerFunc); ok {
		return fn
	}
	return func(c *gin.Context) {
		c.Next()
	}
}
