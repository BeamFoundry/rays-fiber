package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/redirect"
)

// RedirectNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Redirect], to indicates when to skip it this middleware.
//
// Please see the [redirect] middleware configuration documentation.
//
// [redirect]: https://docs.gofiber.io/middleware/redirect/#config
type RedirectNext interface {
	RedirectNext() func(fiber.Ctx) bool
}

// RedirectConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Redirect], to review the [redirect.Config] that will be configured by [Redirect]
//
// Please see the [redirect] middleware configuration documentation.
//
// [redirect]: https://docs.gofiber.io/middleware/redirect/#config
type RedirectConfig interface {
	RedirectConfig(config *redirect.Config)
}

type redirectConfigurationProperties struct {
	rules      []map[string]*string
	statusCode int
}

type redirectConfig struct {
	redirect redirectConfigurationProperties
}

var _, _ = any(&FiberMwConfig[redirectConfig]{}).(Stereotype)

func (this *redirectConfig) Config() redirect.Config {
	config := redirect.ConfigDefault

	ruleList := []redirect.Rule{}
	for _, rules := range this.redirect.rules {
		for k, v := range rules {
			rule := redirect.Rule{From: k, To: *v}
			ruleList = append(ruleList, rule)
		}
	}
	config.RuleList = ruleList
	config.StatusCode = this.redirect.statusCode

	return config
}

// [Redirect] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  Redirect
//	}
//
//	func (this *ApiController) RedirectConfig(config *redirect.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through the [RedirectNext] interface the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
//
// Also, the parent container may implements [RedirectConfig] to receive a pointer to the [redirect.Config]
// that will be used.
//
// Please see those interface definitions and the [redirect] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					redirect:
//						rules:
//							/api/*: null # null values will be ignored.
//						statusCode: 302
//
// [redirect]: https://docs.gofiber.io/middleware/redirect/#config
type Redirect struct {
	config *FiberMwConfig[redirectConfig] `@:"Inject"`
	cfg    *redirect.Config
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *Redirect) PostConstruct(parent any) Error {
	this.cfg = new(this.config.GetConfig(parent).Config())
	if next, ok := parent.(RedirectNext); ok {
		this.cfg.Next = next.RedirectNext()
	}
	if config, ok := parent.(RedirectConfig); ok {
		config.RedirectConfig(this.cfg)
	}
	return nil
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *Redirect) Handler() fiber.Handler {
	return redirect.New(*this.cfg)
}
