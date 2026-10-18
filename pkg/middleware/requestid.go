package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

// RequestIDNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [RequestID], to indicates when to skip it this middleware.
//
// Please see the [requestid] middleware configuration documentation.
//
// [requestid]: https://docs.gofiber.io/middleware/requestid/#config
type RequestIDNext interface {
	RequestIDNext() func(fiber.Ctx) bool
}

// RequestIDNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [RequestID], to indicates when to skip it this middleware.
//
// Please see the [requestid] middleware configuration documentation.
//
// [requestid]: https://docs.gofiber.io/middleware/requestid/#config
type RequestIDGenerator interface {
	RequestIDGenerator() func() string
}

// RequestIDConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [RequestID], to review the [requestid.Config] that will be configured by [RequestID]
//
// Please see the [requestid] middleware configuration documentation.
//
// [requestid]: https://docs.gofiber.io/middleware/requestid/#config
type RequestIDConfig interface {
	RequestIDConfig(config *requestid.Config)
}

type requestIDConfigurationProperties struct {
	header string
}

type requestIDConfig struct {
	requestid requestIDConfigurationProperties
}

var _, _ = any(&FiberMwConfig[requestIDConfig]{}).(Stereotype)

func (this *requestIDConfig) Config() requestid.Config {
	config := requestid.ConfigDefault

	config.Header = this.requestid.header

	return config
}

// [RequestID] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  RequestID
//	}
//
//	func (this *ApiController) RequestIDConfig(config *requestid.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [RequestIDNext] and [RequestIDGenerator].
//
// Also, the parent container may implements [RequestIDConfig] to receive a pointer to the [requestid.Config]
// that will be used.
//
// Please see those interface definitions and the [requestid] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					requestid:
//						header: "X-Request-ID"
//
// [requestid]: https://docs.gofiber.io/middleware/requestid/#config
type RequestID struct {
	config *FiberMwConfig[requestIDConfig] `@:"Inject"`
	cfg    *requestid.Config
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *RequestID) PostConstruct(parent any) Error {
	this.cfg = new(this.config.GetConfig(parent).Config())
	if next, ok := parent.(RequestIDNext); ok {
		this.cfg.Next = next.RequestIDNext()
	}
	if generator, ok := parent.(RequestIDGenerator); ok {
		this.cfg.Generator = generator.RequestIDGenerator()
	}
	if config, ok := parent.(RequestIDConfig); ok {
		config.RequestIDConfig(this.cfg)
	}
	return nil
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *RequestID) Handler() fiber.Handler {
	return requestid.New(*this.cfg)
}
