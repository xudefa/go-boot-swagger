// Package swagger 提供 Swagger API 文档的自动配置。
//
// 创建并注册 Swagger Config Bean 到 IoC 容器中（Bean ID: swaggerConfig）。
package swagger

import (
	swaggercore "github.com/xudefa/go-boot-swagger"

	"github.com/xudefa/go-boot/core"
)

const (
	BeanSwaggerConfig = "swaggerConfig"
)

// AutoConfig 自动配置Swagger
func AutoConfig(container core.Container) error {
	config := swaggercore.DefaultConfig()
	return container.Register(BeanSwaggerConfig, core.Bean(config), core.Singleton())
}
