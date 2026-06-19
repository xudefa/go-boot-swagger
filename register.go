package swagger

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// RegisterAuto 自动适配并注册Swagger路由
func RegisterAuto(engine interface{}, opts ...Option) error {
	adapter, err := globalRegistry.GetAdapterByEngine(engine)
	if err != nil {
		return fmt.Errorf("failed to get adapter: %w", err)
	}

	config := DefaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	return adapter.Register(engine, config.Path, opts...)
}

// RegisterWithGroupAuto 自动适配并在路由组中注册Swagger
func RegisterWithGroupAuto(group interface{}, opts ...Option) error {
	adapter, err := globalRegistry.GetAdapterByGroup(group)
	if err != nil {
		return fmt.Errorf("failed to get adapter: %w", err)
	}

	config := DefaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	return adapter.RegisterWithGroup(group, config.Path, opts...)
}

// MiddlewareAuto 自动适配并返回Swagger中间件
func MiddlewareAuto(engineType EngineType, opts ...Option) (interface{}, error) {
	adapter, err := globalRegistry.GetAdapter(engineType)
	if err != nil {
		return nil, fmt.Errorf("failed to get adapter: %w", err)
	}

	return adapter.Middleware(opts...), nil
}

// HealthCheckAuto 自动适配的健康检查端点
func HealthCheckAuto(c interface{}) {
	for _, adapter := range globalRegistry.adapters {
		adapter.HealthCheck(c)
	}
}

// CORSForSwaggerAuto 自动适配的CORS中间件
func CORSForSwaggerAuto(engineType EngineType) (interface{}, error) {
	adapter, err := globalRegistry.GetAdapter(engineType)
	if err != nil {
		return nil, fmt.Errorf("failed to get adapter: %w", err)
	}

	return adapter.CORSForSwagger(), nil
}

// RegisterForGin 为Gin框架注册Swagger（向后兼容）
func RegisterForGin(router *gin.Engine, opts ...Option) {
	adapter := NewGinAdapter()
	config := DefaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	if !config.Enabled {
		return
	}

	middleware := adapter.Middleware(opts...)
	if handler, ok := middleware.(gin.HandlerFunc); ok {
		router.GET(config.Path, handler)
	}
}

// RegisterWithGroupForGin 为Gin路由组注册Swagger（向后兼容）
func RegisterWithGroupForGin(group *gin.RouterGroup, opts ...Option) {
	adapter := NewGinAdapter()
	config := DefaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	if !config.Enabled {
		return
	}

	middleware := adapter.Middleware(opts...)
	if handler, ok := middleware.(gin.HandlerFunc); ok {
		group.GET(config.Path, handler)
	}
}

// MiddlewareForGin 为Gin框架返回Swagger中间件（向后兼容）
func MiddlewareForGin(opts ...Option) interface{} {
	return NewGinAdapter().Middleware(opts...)
}

// RegisterForHertz 为Hertz框架注册Swagger
func RegisterForHertz(router interface{}, opts ...Option) error {
	return RegisterAuto(router, opts...)
}

// RegisterWithGroupForHertz 为Hertz路由组注册Swagger
func RegisterWithGroupForHertz(group interface{}, opts ...Option) error {
	return RegisterWithGroupAuto(group, opts...)
}

// MiddlewareForHertz 为Hertz框架返回Swagger中间件
func MiddlewareForHertz(opts ...Option) (interface{}, error) {
	return MiddlewareAuto(EngineTypeHertz, opts...)
}
