package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/etag"
)

// EtagNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Etag], to indicates when to skip it this middleware.
//
// Please see the [etag] middleware configuration documentation.
//
// [etag]: https://docs.gofiber.io/middleware/etag/#config
type EtagNext interface {
	EtagNext() func(fiber.Ctx) bool
}

// EtagConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Etag], to review the [etag.Config] that will be configured by [Etag]
//
// Please see the [etag] middleware configuration documentation.
//
// [etag]: https://docs.gofiber.io/middleware/etag/#config
type EtagConfig interface {
	EtagConfig(config *etag.Config)
}

type etagConfigurationProperties struct {
	week bool
}

type etagConfig struct {
	etag etagConfigurationProperties
}

var _, _ = any(&FiberMwConfig[etagConfig]{}).(Stereotype)

func (this *etagConfig) Config() etag.Config {
	config := etag.ConfigDefault
	config.Weak = this.etag.week

	return config
}

// [Etag] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  Etag
//	}
//
//	func (this *ApiController) EtagConfig(config *etag.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// The only interface is [EtagNext].
//
// Also, the parent container may implements [EtagConfig] to receive a pointer to the [etag.Config]
// that will be used.
//
// Please see those interface definitions and the [etag] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					etag:
//						weak: false
//
// [etag]: https://docs.gofiber.io/middleware/etag/#config
type Etag struct {
	config *FiberMwConfig[etagConfig] `@:"Inject"`
	cfg    *etag.Config
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *Etag) PostConstruct(parent any) Error {
	this.cfg = new(this.config.GetConfig(parent).Config())
	if next, ok := parent.(EtagNext); ok {
		this.cfg.Next = next.EtagNext()
	}
	if config, ok := parent.(EtagConfig); ok {
		config.EtagConfig(this.cfg)
	}
	return nil
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *Etag) Handler() fiber.Handler {
	return etag.New(*this.cfg)
}
