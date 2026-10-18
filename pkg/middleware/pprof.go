package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/pprof"
)

// PprofNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Pprof], to indicates when to skip it this middleware.
//
// Please see the [pprof] middleware configuration documentation.
//
// [pprof]: https://docs.gofiber.io/middleware/pprof/#config
type PprofNext interface {
	PprofNext() func(fiber.Ctx) bool
}

// PprofConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Pprof], to review the [pprof.Config] that will be configured by [Pprof]
//
// Please see the [pprof] middleware configuration documentation.
//
// [pprof]: https://docs.gofiber.io/middleware/pprof/#config
type PprofConfig interface {
	PprofConfig(config *pprof.Config)
}

type pprofConfigurationProperties struct {
	prefix string
}

var _, _ = any(&pprofConfigurationProperties{}).(Stereotype)

type pprofConfig struct {
	pprof pprofConfigurationProperties
}

var _, _ = any(&FiberMwConfig[pprofConfig]{}).(Stereotype)

func (this *pprofConfig) Config() pprof.Config {
	config := pprof.ConfigDefault

	config.Prefix = this.pprof.prefix

	return config
}

// [Pprof] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  Pprof
//	}
//
//	func (this *ApiController) PprofConfig(config *pprof.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through the interface [PprofNext] the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
//
// Also, the parent container may implements [PprofConfig] to receive a pointer to the [pprof.Config]
// that will be used.
//
// Please see those interface definitions and the [pprof] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					pprof:
//						prefix: ""
//
// [pprof]: https://docs.gofiber.io/middleware/pprof/#config
type Pprof struct {
	config *FiberMwConfig[pprofConfig] `@:"Inject"`
	cfg    *pprof.Config
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *Pprof) PostConstruct(parent any) Error {
	this.cfg = new(this.config.GetConfig(parent).Config())
	if next, ok := parent.(PprofNext); ok {
		this.cfg.Next = next.PprofNext()
	}
	if config, ok := parent.(PprofConfig); ok {
		config.PprofConfig(this.cfg)
	}
	return nil
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *Pprof) Handler() fiber.Handler {
	return pprof.New(*this.cfg)
}
