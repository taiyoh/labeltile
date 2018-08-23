package labeltile

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/taiyoh/labeltile/app"
	"github.com/taiyoh/labeltile/app/infra"
)

// Labeltile is main interface of this application
type Labeltile struct {
	container app.Container
	engine    *gin.Engine
}

// ContainerInjector is alias for infra module injection using Functional Option Pattern
type ContainerInjector func(app.Container, *Conf)

// RouteInjector is alias for binding endpoint and controller implementation using Functional Option Pattern
type RouteInjector func(*gin.Engine, app.Container)

// NewLabeltile returns Labeltile object
func NewLabeltile(conf *Conf, injectFns ...ContainerInjector) *Labeltile {
	container := infra.NewContainer()
	for _, fn := range injectFns {
		fn(container, conf)
	}
	router := gin.Default()
	router.LoadHTMLGlob(filepath.Join(conf.Server.Template, "*.tmpl"))
	router.Use(SetupUserTokenMiddleware(container))

	return &Labeltile{
		container: container,
		engine:    router,
	}
}

// SetupRoutes provides filling routing in gin framework
func (l *Labeltile) SetupRoutes(injectFns ...RouteInjector) {
	for _, fn := range injectFns {
		fn(l.engine, l.container)
	}
}
