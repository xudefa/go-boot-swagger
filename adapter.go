package swagger

import (
	"errors"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/gin-gonic/gin"
)

var (
	// ErrInvalidEngine 无效的引擎类型
	ErrInvalidEngine = errors.New("invalid engine type")
	// ErrInvalidGroup 无效的路由组类型
	ErrInvalidGroup = errors.New("invalid router group type")
	// ErrUnsupportedFramework 不支持的框架
	ErrUnsupportedFramework = errors.New("unsupported framework")
)

// EngineType 引擎类型枚举
type EngineType int

const (
	// EngineTypeGin Gin框架
	EngineTypeGin EngineType = iota
	// EngineTypeHertz Hertz框架
	EngineTypeHertz
	// EngineTypeEcho Echo框架
	EngineTypeEcho
	// EngineTypeFiber Fiber框架
	EngineTypeFiber
	// EngineTypeChi Chi框架
	EngineTypeChi
)

// FrameworkAdapter 框架适配器接口
type FrameworkAdapter interface {
	// Register 注册Swagger路由到引擎
	Register(engine interface{}, path string, opts ...Option) error

	// RegisterWithGroup 在路由组中注册Swagger
	RegisterWithGroup(group interface{}, path string, opts ...Option) error

	// Middleware 返回Swagger中间件
	Middleware(opts ...Option) interface{}

	// HealthCheck 健康检查端点
	HealthCheck(c interface{})

	// CORSForSwagger 为Swagger UI提供CORS支持
	CORSForSwagger() interface{}

	// GetEngineType 返回引擎类型
	GetEngineType() EngineType
}

// AdapterRegistry 适配器注册表
type AdapterRegistry struct {
	adapters map[EngineType]FrameworkAdapter // 适配器映射
}

// NewAdapterRegistry 创建适配器注册表
func NewAdapterRegistry() *AdapterRegistry {
	registry := &AdapterRegistry{
		adapters: make(map[EngineType]FrameworkAdapter),
	}

	// 注册默认适配器
	registry.RegisterAdapter(EngineTypeGin, NewGinAdapter())
	registry.RegisterAdapter(EngineTypeHertz, NewHertzAdapter())

	return registry
}

// RegisterAdapter 注册适配器
func (r *AdapterRegistry) RegisterAdapter(engineType EngineType, adapter FrameworkAdapter) {
	r.adapters[engineType] = adapter
}

// GetAdapter 获取适配器
func (r *AdapterRegistry) GetAdapter(engineType EngineType) (FrameworkAdapter, error) {
	adapter, ok := r.adapters[engineType]
	if !ok {
		return nil, ErrUnsupportedFramework
	}
	return adapter, nil
}

// GetAdapterByEngine 根据引擎实例获取适配器
func (r *AdapterRegistry) GetAdapterByEngine(engine interface{}) (FrameworkAdapter, error) {
	switch engine.(type) {
	case *gin.Engine:
		return r.GetAdapter(EngineTypeGin)
	case *server.Hertz:
		return r.GetAdapter(EngineTypeHertz)
	default:
		return nil, ErrUnsupportedFramework
	}
}

// GetAdapterByGroup 根据路由组实例获取适配器
func (r *AdapterRegistry) GetAdapterByGroup(group interface{}) (FrameworkAdapter, error) {
	switch group.(type) {
	case *gin.RouterGroup:
		return r.GetAdapter(EngineTypeGin)
	case *route.RouterGroup:
		return r.GetAdapter(EngineTypeHertz)
	default:
		return nil, ErrUnsupportedFramework
	}
}

// 全局适配器注册表
var globalRegistry = NewAdapterRegistry()

// RegisterAdapter 注册全局适配器
func RegisterAdapter(engineType EngineType, adapter FrameworkAdapter) {
	globalRegistry.RegisterAdapter(engineType, adapter)
}

// GetAdapter 获取全局适配器
func GetAdapter(engineType EngineType) (FrameworkAdapter, error) {
	return globalRegistry.GetAdapter(engineType)
}
