package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

// RecoverNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Recover], to indicates when to skip it this middleware.
//
// Please see the [recover] middleware configuration documentation.
//
// [recover]: https://docs.gofiber.io/middleware/recover/#config
type RecoverNext interface {
	RecoverNext() func(fiber.Ctx) bool
}

// RecoverPanicHandler is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Recover], to indicates a customize the error returned from a recovered panic.
//
// Please see the [recover] middleware configuration documentation.
//
// [recover]: https://docs.gofiber.io/middleware/recover/#config
type RecoverPanicHandler interface {
	RecoverPanicHandler() func(c fiber.Ctx, e any) error
}

// RecoverStackTraceHandler is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Recover], to indicates how to handle the captured stack trace when enabled.
//
// Please see the [recover] middleware configuration documentation.
//
// [recover]: https://docs.gofiber.io/middleware/recover/#config
type RecoverStackTraceHandler interface {
	RecoverStackTraceHandler() func(c fiber.Ctx, e any)
}

// RecoverConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Recover], to review the [recover.Config] that will be configured by [Recover]
//
// Please see the [recover] middleware configuration documentation.
//
// [recover]: https://docs.gofiber.io/middleware/recover/#config
type RecoverConfig interface {
	RecoverConfig(config *recover.Config)
}

type recoverConfigurationProperties struct {
	enableStackTrace bool
}

type recoverConfig struct {
	recovering recoverConfigurationProperties
}

var _, _ = any(&FiberMwConfig[recoverConfig]{}).(Stereotype)

func (this *recoverConfig) Config() recover.Config {
	config := recover.ConfigDefault
	config.EnableStackTrace = this.recovering.enableStackTrace

	return config
}

// [Recover] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  Recover
//	}
//
//	func (this *ApiController) RecoverConfig(config *recover.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [RecoverNext] [RecoverPanicHandler] and [RecoverStackTraceHandler].
//
// Also, the parent container may implements [recoverConfig] to receive a pointer to the [recover.Config]
// that will be used.
//
// Please see those interface definitions and the [recover] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					recover:
//						enableStackTrace: false
//
// [recover]: https://docs.gofiber.io/middleware/recover/#config
type Recover struct {
	config *FiberMwConfig[recoverConfig] `@:"Inject"`
	cfg    *recover.Config
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *Recover) PostConstruct(parent any) Error {
	this.cfg = new(this.config.GetConfig(parent).Config())
	if next, ok := (parent).(RecoverNext); ok {
		this.cfg.Next = next.RecoverNext()
	}
	if panicHandler, ok := (parent).(RecoverPanicHandler); ok {
		this.cfg.PanicHandler = panicHandler.RecoverPanicHandler()
	}
	if stackTraceHandler, ok := (parent).(RecoverStackTraceHandler); ok {
		this.cfg.StackTraceHandler = stackTraceHandler.RecoverStackTraceHandler()
	}
	if config, ok := (parent).(RecoverConfig); ok {
		config.RecoverConfig(this.cfg)
	}
	return nil
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *Recover) Handler() fiber.Handler {
	return recover.New(*this.cfg)
}
