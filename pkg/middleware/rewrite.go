package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/rewrite"
)

// RewriteNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Rewrite], to indicates when to skip it this middleware.
//
// Please see the [rewrite] middleware configuration documentation.
//
// [rewrite]: https://docs.gofiber.io/middleware/rewrite/#config
type RewriteNext interface {
	RewriteNext() func(fiber.Ctx) bool
}

// RewriteConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Rewrite], to review the [rewrite.Config] that will be configured by [Rewrite]
//
// Please see the [rewrite] middleware configuration documentation.
//
// [rewrite]: https://docs.gofiber.io/middleware/rewrite/#config
type RewriteConfig interface {
	RewriteConfig(config *rewrite.Config)
}

type rewriteConfigurationProperties struct {
	rules map[string]*string
}

type rewriteConfig struct {
	rewrite rewriteConfigurationProperties
	//
}

var _, _ = any(&FiberMwConfig[rewriteConfig]{}).(Stereotype)

func (this *rewriteConfig) Config() rewrite.Config {
	config := rewrite.Config{}

	rules := map[string]string{}
	for k, v := range this.rewrite.rules {
		if v != nil {
			rules[k] = *v
		}
	}
	config.Rules = rules

	return config
}

// [Rewrite] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  Rewrite
//	}
//
//	func (this *ApiController) RewriteConfig(config *rewrite.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through the [RewriteNext] interface the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
//
// Also, the parent container may implements [RewriteConfig] to receive a pointer to the [rewrite.Config]
// that will be used.
//
// Please see those interface definitions and the [rewrite] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					rewrite:
//	          rules:
//	            /api/*: null # null values will be ignored.
//
// [rewrite]: https://docs.gofiber.io/middleware/rewrite/#config
type Rewrite struct {
	config *FiberMwConfig[rewriteConfig] `@:"Inject"`
	cfg    *rewrite.Config
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *Rewrite) PostConstruct(parent any) Error {
	this.cfg = new(this.config.GetConfig(parent).Config())
	if next, ok := parent.(RewriteNext); ok {
		this.cfg.Next = next.RewriteNext()
	}
	if config, ok := parent.(RewriteConfig); ok {
		config.RewriteConfig(this.cfg)
	}
	return nil
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *Rewrite) Handler() fiber.Handler {
	return rewrite.New(*this.cfg)
}
