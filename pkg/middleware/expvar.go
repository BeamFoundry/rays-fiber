package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/lang"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/expvar"
)

// ExpVarNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [ExpVar], to indicates when to skip it this middleware.
//
// Please see the [expvar] middleware configuration documentation.
//
// [expvar]: https://docs.gofiber.io/middleware/expvar/#config
type ExpVarNext interface {
	ExpVarNext() func(fiber.Ctx) bool
}

// ExpVarConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [ExpVar], to review the [expvar.Config] that will be configured by [ExpVar]
//
// Please see the [expvar] middleware configuration documentation.
//
// [expvar]: https://docs.gofiber.io/middleware/expvar/#config
type ExpVarConfig interface {
	ExpVarConfig(config *expvar.Config)
}

// [ExpVar] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  ExpVar
//	}
//
//	func (this *ApiController) ExpVarConfig(config *expvar.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through the interface [ExpVarNext] the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
//
// Also, the parent container may implements [ExpVarConfig] to receive a pointer to the [expvar.Config]
// that will be used.
//
// Please see those interface definitions and the [expvar] middleware documentation.
//
// [expvar]: https://docs.gofiber.io/middleware/expvar/#config
type ExpVar struct {
	cfg *expvar.Config
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *ExpVar) PostConstruct(parent any) Error {
	return Try(func() {
		this.cfg = new(expvar.ConfigDefault)
		if next, ok := parent.(ExpVarNext); ok {
			this.cfg.Next = next.ExpVarNext()
		}
	})
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *ExpVar) Handler() fiber.Handler {
	return expvar.New(*this.cfg)
}
