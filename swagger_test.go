// /Users/xudefa/workspace/go-boot/swagger/swagger_test.go 中添加更多测试
// 我们将扩展现有的测试文件
package swagger_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/xudefa/go-boot-swagger"
)

func TestDefaultConfig(t *testing.T) {
	config := swagger.DefaultConfig()
	assert.NotNil(t, config)
	assert.True(t, config.Enabled)
	assert.Equal(t, "/swagger/*any", config.Path)
	assert.Equal(t, "API Documentation", config.Title)
	assert.Equal(t, "localhost:8080", config.Host)
}

func TestWithEnabled(t *testing.T) {
	config := swagger.DefaultConfig()
	opt := swagger.WithEnabled(false)
	opt(config)
	assert.False(t, config.Enabled)
}

func TestWithPath(t *testing.T) {
	config := swagger.DefaultConfig()
	opt := swagger.WithPath("/docs/*any")
	opt(config)
	assert.Equal(t, "/docs/*any", config.Path)
}

func TestWithTitle(t *testing.T) {
	config := swagger.DefaultConfig()
	opt := swagger.WithTitle("My API")
	opt(config)
	assert.Equal(t, "My API", config.Title)
}

func TestWithDescription(t *testing.T) {
	config := swagger.DefaultConfig()
	opt := swagger.WithDescription("My API Description")
	opt(config)
	assert.Equal(t, "My API Description", config.Description)
}

func TestWithVersion(t *testing.T) {
	config := swagger.DefaultConfig()
	opt := swagger.WithVersion("2.0")
	opt(config)
	assert.Equal(t, "2.0", config.Version)
}

func TestWithHost(t *testing.T) {
	config := swagger.DefaultConfig()
	opt := swagger.WithHost("api.example.com")
	opt(config)
	assert.Equal(t, "api.example.com", config.Host)
}

func TestWithBasePath(t *testing.T) {
	config := swagger.DefaultConfig()
	opt := swagger.WithBasePath("/api/v1")
	opt(config)
	assert.Equal(t, "/api/v1", config.BasePath)
}

func TestWithSchemes(t *testing.T) {
	config := swagger.DefaultConfig()
	opt := swagger.WithSchemes("https", "wss")
	opt(config)
	assert.Equal(t, []string{"https", "wss"}, config.Schemes)
}

func TestWithSecurityDefinitions(t *testing.T) {
	defs := map[string]swagger.SecurityDefinition{
		"ApiKeyAuth": {
			Type:        "apiKey",
			Name:        "X-API-Key",
			In:          "header",
			Description: "API Key Authentication",
		},
	}

	config := swagger.DefaultConfig()
	opt := swagger.WithSecurityDefinitions(defs)
	opt(config)
	assert.Equal(t, defs, config.SecurityDefinitions)
}

func TestWithContactInfo(t *testing.T) {
	config := swagger.DefaultConfig()
	opt := swagger.WithContact("John Doe", "john@example.com", "https://example.com")
	opt(config)

	assert.Equal(t, "John Doe", config.ContactName)
	assert.Equal(t, "john@example.com", config.ContactEmail)
	assert.Equal(t, "https://example.com", config.ContactURL)
}

func TestWithLicenseInfo(t *testing.T) {
	config := swagger.DefaultConfig()
	opt := swagger.WithLicense("MIT", "https://opensource.org/licenses/MIT")
	opt(config)

	assert.Equal(t, "MIT", config.LicenseName)
	assert.Equal(t, "https://opensource.org/licenses/MIT", config.LicenseURL)
}

func TestWithSecurityDefinitionsBatch(t *testing.T) {
	defs := map[string]swagger.SecurityDefinition{
		"ApiKeyAuth": {
			Type:        "apiKey",
			Name:        "X-API-Key",
			In:          "header",
			Description: "API Key Authentication",
		},
		"OAuth2": {
			Type:             "oauth2",
			AuthorizationURL: "https://example.com/oauth/authorize",
			TokenURL:         "https://example.com/oauth/token",
			Scopes: map[string]string{
				"read":  "Grants read access",
				"write": "Grants write access",
			},
		},
	}

	config := swagger.DefaultConfig()
	opt := swagger.WithSecurityDefinitions(defs)
	opt(config)

	assert.Equal(t, defs, config.SecurityDefinitions)
}

func TestUpdateConfig(t *testing.T) {
	config := swagger.DefaultConfig()
	newConfig := &swagger.Config{
		Enabled:      false,
		Title:        "Updated Title",
		Description:  "Updated Description",
		Version:      "3.0",
		Host:         "newhost.com",
		BasePath:     "/new/api",
		Schemes:      []string{"https"},
		ContactName:  "Jane Doe",
		ContactEmail: "jane@example.com",
		ContactURL:   "https://jane.com",
		LicenseName:  "Apache 2.0",
		LicenseURL:   "https://www.apache.org/licenses/LICENSE-2.0",
	}

	config.UpdateConfig(newConfig)

	assert.Equal(t, false, config.Enabled)
	assert.Equal(t, "Updated Title", config.Title)
	assert.Equal(t, "Updated Description", config.Description)
	assert.Equal(t, "3.0", config.Version)
	assert.Equal(t, "newhost.com", config.Host)
	assert.Equal(t, "/new/api", config.BasePath)
	assert.Equal(t, []string{"https"}, config.Schemes)
	assert.Equal(t, "Jane Doe", config.ContactName)
	assert.Equal(t, "jane@example.com", config.ContactEmail)
	assert.Equal(t, "https://jane.com", config.ContactURL)
	assert.Equal(t, "Apache 2.0", config.LicenseName)
	assert.Equal(t, "https://www.apache.org/licenses/LICENSE-2.0", config.LicenseURL)
}

func TestMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	middleware := swagger.Middleware(swagger.WithEnabled(true))
	router.Use(middleware)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMiddlewareDisabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	middleware := swagger.Middleware(swagger.WithEnabled(false))
	router.Use(middleware)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	swagger.Register(router, swagger.WithEnabled(true))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRegisterDisabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	swagger.Register(router, swagger.WithEnabled(false))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
