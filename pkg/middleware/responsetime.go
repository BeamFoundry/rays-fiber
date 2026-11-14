package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/responsetime"
)

// ResponseTimeNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [ResponseTime], to indicates when to skip it this middleware.
//
// Please see the [responsetime] middleware configuration documentation.
//
// [responsetime]: https://docs.gofiber.io/middleware/responsetime/#config
type ResponseTimeNext interface {
	ResponseTimeNext() func(fiber.Ctx) bool
}

// ResponseTimeConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [ResponseTime], to review the [responsetime.Config] that will be configured by [responsetime]
//
// Please see the [responsetime] middleware configuration documentation.
//
// [responsetime]: https://docs.gofiber.io/middleware/responsetime/#config
type ResponseTimeConfig interface {
	ResponseTimeConfig(config *responsetime.Config)
}

type responsetimeConfigurationProperties struct {
	header string
}

type responsetimeConfig struct {
	responsetime responsetimeConfigurationProperties
}

var _, _ = any(&FiberMwConfig[responsetimeConfig]{}).(Stereotype)

func (this *responsetimeConfig) Config() responsetime.Config {
	config := responsetime.ConfigDefault

	config.Header = this.responsetime.header

	return config
}

// [ResponseTime] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//		FiberController
//		RequestMapping `Path:"/api"`
//		ResponseTime
//	}
//
//	func (this *ApiController) ResponseTimeConfig(config *responsetime.Config) {
//		// changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [ResponseTimeNext], [], [] and [].
//
// Also, the parent container may implements [ResponseTimeConfig] to receive a pointer to the [responsetime.Config]
// that will be used.
//
// Please see those interface definitions and the [responsetime] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					responsetime:
//						header: 5m
//
// [responsetime]: https://docs.gofiber.io/middleware/responsetime/#config
type ResponseTime struct {
	config *FiberMwConfig[responsetimeConfig] `@:"Inject"`
	cfg    *responsetime.Config
}

func (this *ResponseTime) PostConstruct(parent any) Error {
	cfg := this.config.GetConfig(parent)
	this.cfg = new(cfg.Config())
	if next, ok := parent.(ResponseTimeNext); ok {
		this.cfg.Next = next.ResponseTimeNext()
	}
	if config, ok := parent.(ResponseTimeConfig); ok {
		config.ResponseTimeConfig(this.cfg)
	}
	return nil
}

func (this *ResponseTime) Handler() fiber.Handler {
	return responsetime.New(*this.cfg)
}
