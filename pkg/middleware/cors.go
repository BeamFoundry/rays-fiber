package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

// CorsNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Cors], to indicates when to skip it this middleware.
//
// Please see the [cors] middleware configuration documentation.
//
// [cors]: https://docs.gofiber.io/middleware/cors/#config
type CorsNext interface {
	CorsNext() func(fiber.Ctx) bool
}

// CorsAllowOriginsFunc is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Cors], to indicates a function that dynamically determines whether
// to allow a request based on its origin.
//
// Please see the [cors] middleware configuration documentation.
//
// [cors]: https://docs.gofiber.io/middleware/cors/#config
type CorsAllowOriginsFunc interface {
	CorsAllowOriginsFunc() func(origin string) bool
}

// CorsConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Cors], to review the [cors.Config] that will be configured by [Cors]
//
// Please see the [cors] middleware configuration documentation.
//
// [cors]: https://docs.gofiber.io/middleware/cors/#config
type CorsConfig interface {
	CorsConfig(config *cors.Config)
}

type corsConfigurationProperties struct {
	allowOrigins          []string
	allowMethods          []string
	allowHeaders          []string
	exposeHeaders         []string
	maxAge                int
	allowCredentials      bool
	allowPrivateNetwork   bool
	disableValueRedaction bool
}

type corsConfig struct {
	cors corsConfigurationProperties
}

var _, _ = any(&FiberMwConfig[corsConfig]{}).(Stereotype)

func (this *corsConfig) Config() cors.Config {
	config := cors.ConfigDefault

	config.AllowOrigins = this.cors.allowOrigins
	config.AllowMethods = this.cors.allowMethods
	config.AllowHeaders = this.cors.allowHeaders
	config.ExposeHeaders = this.cors.exposeHeaders
	config.MaxAge = this.cors.maxAge
	config.AllowCredentials = this.cors.allowCredentials
	config.AllowPrivateNetwork = this.cors.allowPrivateNetwork
	config.DisableValueRedaction = this.cors.disableValueRedaction

	return config
}

// [Cors] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//		FiberController
//		RequestMapping `Path:"/api"`
//		Cors
//	}
//
//	func (this *ApiController) CorsConfig(config *cors.Config) {
//		// changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [CorsNext], and [CorsAllowOriginsFunc].
//
// Also, the parent container may implements [CorsConfig] to receive a pointer to the [cors.Config]
// that will be used.
//
// Please see those interface definitions and the [cors] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					cors:
//						allowOrigins: ["*"]
//						allowMethods: ["GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"]
//						allowHeaders: []
//						allowCredentials: false
//						allowPrivateNetwork: false
//						disableValueRedaction: false
//						exposeHeaders: []
//						maxAge: 0
//
// [cors]: https://docs.gofiber.io/middleware/cors/#config
type Cors struct {
	config *FiberMwConfig[corsConfig] `@:"Inject"`
	cfg    *cors.Config
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *Cors) PostConstruct(parent any) Error {
	this.cfg = new(this.config.GetConfig(parent).Config())
	if next, ok := (parent).(CorsNext); ok {
		this.cfg.Next = next.CorsNext()
	}
	if allowOriginsFunc, ok := (parent).(CorsAllowOriginsFunc); ok {
		this.cfg.AllowOriginsFunc = allowOriginsFunc.CorsAllowOriginsFunc()
	}
	if config, ok := (parent).(CorsConfig); ok {
		config.CorsConfig(this.cfg)
	}
	return nil
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *Cors) Handler() fiber.Handler {
	return cors.New(*this.cfg)
}
