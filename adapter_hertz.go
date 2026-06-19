package swagger

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/route"
	swaggerFiles "github.com/swaggo/files"
)

// HertzAdapter Hertz框架适配器
type HertzAdapter struct{}

// NewHertzAdapter 创建Hertz适配器
func NewHertzAdapter() *HertzAdapter {
	return &HertzAdapter{}
}

// Register 注册Swagger路由到Hertz引擎
func (a *HertzAdapter) Register(engine interface{}, path string, opts ...Option) error {
	router, ok := engine.(*server.Hertz)
	if !ok {
		return ErrInvalidEngine
	}

	config := DefaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	if !config.Enabled {
		return nil
	}

	middleware := a.Middleware(opts...)
	if handler, ok := middleware.(app.HandlerFunc); ok {
		router.GET(path, handler)
	}
	return nil
}

// RegisterWithGroup 在路由组中注册Swagger
func (a *HertzAdapter) RegisterWithGroup(group interface{}, path string, opts ...Option) error {
	routerGroup, ok := group.(*route.RouterGroup)
	if !ok {
		return ErrInvalidGroup
	}

	config := DefaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	if !config.Enabled {
		return nil
	}

	middleware := a.Middleware(opts...)
	if handler, ok := middleware.(app.HandlerFunc); ok {
		routerGroup.GET(path, handler)
	}
	return nil
}

// Middleware 返回Swagger中间件
func (a *HertzAdapter) Middleware(opts ...Option) interface{} {
	config := DefaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	if !config.Enabled {
		return func(ctx context.Context, c *app.RequestContext) {
			c.Next(ctx)
		}
	}

	return func(ctx context.Context, c *app.RequestContext) {
		requestPath := string(c.Request.URI().Path())

		if strings.HasPrefix(requestPath, "/swagger/") {
			filePath := strings.TrimPrefix(requestPath, "/swagger/")
			if filePath == "" || filePath == "/" {
				filePath = "index.html"
			}

			file, err := swaggerFiles.HTTP.Open(filePath)
			if err != nil {
				c.SetStatusCode(http.StatusNotFound)
				_, _ = c.WriteString("404 Not Found")
				return
			}
			defer func() { _ = file.Close() }()

			stat, err := file.Stat()
			if err != nil {
				c.SetStatusCode(http.StatusInternalServerError)
				_, _ = c.WriteString("500 Internal Server Error")
				return
			}

			contentType := getContentType(filePath)
			c.SetContentType(contentType)

			if stat.IsDir() {
				c.SetStatusCode(http.StatusForbidden)
				_, _ = c.WriteString("403 Forbidden")
				return
			}

			buf := new(bytes.Buffer)
			if _, err := io.Copy(buf, file); err != nil {
				c.SetStatusCode(http.StatusInternalServerError)
				_, _ = c.WriteString("500 Internal Server Error")
				return
			}

			_, _ = c.Write(buf.Bytes())
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}

// HealthCheck 健康检查端点
func (a *HertzAdapter) HealthCheck(c interface{}) {
	ctx, ok := c.(*app.RequestContext)
	if !ok {
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

// CORSForSwagger 为Swagger UI提供CORS支持
func (a *HertzAdapter) CORSForSwagger() interface{} {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if string(c.Method()) == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next(ctx)
	}
}

// GetEngineType 返回引擎类型
func (a *HertzAdapter) GetEngineType() EngineType {
	return EngineTypeHertz
}

// getContentType 根据文件扩展名获取Content-Type
func getContentType(filePath string) string {
	if strings.HasSuffix(filePath, ".html") {
		return "text/html; charset=utf-8"
	}
	if strings.HasSuffix(filePath, ".css") {
		return "text/css; charset=utf-8"
	}
	if strings.HasSuffix(filePath, ".js") {
		return "application/javascript; charset=utf-8"
	}
	if strings.HasSuffix(filePath, ".json") {
		return "application/json; charset=utf-8"
	}
	if strings.HasSuffix(filePath, ".png") {
		return "image/png"
	}
	if strings.HasSuffix(filePath, ".jpg") || strings.HasSuffix(filePath, ".jpeg") {
		return "image/jpeg"
	}
	if strings.HasSuffix(filePath, ".gif") {
		return "image/gif"
	}
	if strings.HasSuffix(filePath, ".svg") {
		return "image/svg+xml"
	}
	if strings.HasSuffix(filePath, ".ico") {
		return "image/x-icon"
	}
	if strings.HasSuffix(filePath, ".woff") {
		return "font/woff"
	}
	if strings.HasSuffix(filePath, ".woff2") {
		return "font/woff2"
	}
	if strings.HasSuffix(filePath, ".ttf") {
		return "font/ttf"
	}
	if strings.HasSuffix(filePath, ".eot") {
		return "application/vnd.ms-fontobject"
	}
	return "application/octet-stream"
}
