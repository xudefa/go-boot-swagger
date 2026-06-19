package swagger

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// GinAdapter Gin框架适配器
type GinAdapter struct{}

// NewGinAdapter 创建Gin适配器
func NewGinAdapter() *GinAdapter {
	return &GinAdapter{}
}

// Register 注册Swagger路由到Gin引擎
func (a *GinAdapter) Register(engine interface{}, path string, opts ...Option) error {
	router, ok := engine.(*gin.Engine)
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
	if handler, ok := middleware.(gin.HandlerFunc); ok {
		router.GET(path, handler)
	}
	return nil
}

// RegisterWithGroup 在路由组中注册Swagger
func (a *GinAdapter) RegisterWithGroup(group interface{}, path string, opts ...Option) error {
	routerGroup, ok := group.(*gin.RouterGroup)
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
	if handler, ok := middleware.(gin.HandlerFunc); ok {
		routerGroup.GET(path, handler)
	}
	return nil
}

// Middleware 返回Swagger中间件
func (a *GinAdapter) Middleware(opts ...Option) interface{} {
	config := DefaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	if !config.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL(config.Path),
		ginSwagger.DefaultModelsExpandDepth(-1),
	)
}

// HealthCheck 健康检查端点
func (a *GinAdapter) HealthCheck(c interface{}) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// CORSForSwagger 为Swagger UI提供CORS支持
func (a *GinAdapter) CORSForSwagger() interface{} {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// GetEngineType 返回引擎类型
func (a *GinAdapter) GetEngineType() EngineType {
	return EngineTypeGin
}
