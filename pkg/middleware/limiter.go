package middleware

import (
	"time"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

// LimiterNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Limiter], to indicates when to skip it this middleware.
//
// Please see the [limiter] middleware configuration documentation.
//
// [limiter]: https://docs.gofiber.io/middleware/limiter/#config
type LimiterNext interface {
	LimiterNext() func(fiber.Ctx) bool
}

// LimiterMaxFunc is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Limiter], to indicates a function that calculates the maximum number of recent
// connections within Expiration seconds before sending a 429 response.
//
// Please see the [limiter] middleware configuration documentation.
//
// [limiter]: https://docs.gofiber.io/middleware/limiter/#config
type LimiterMaxFunc interface {
	LimiterMaxFunc() func(c fiber.Ctx) int
}

// LimiterKeyGenerator is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Limiter], to indicates a function to generate custom keys; uses Ctx.IP() by default.
//
// Please see the [limiter] middleware configuration documentation.
//
// [limiter]: https://docs.gofiber.io/middleware/limiter/#config
type LimiterKeyGenerator interface {
	LimiterKeyGenerator() func(fiber.Ctx) string
}

// LimiterLimitReached is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Limiter], to indicates a function called when a request exceeds the limit.
//
// Please see the [limiter] middleware configuration documentation.
//
// [limiter]: https://docs.gofiber.io/middleware/limiter/#config
type LimiterLimitReached interface {
	LimiterLimitReached() fiber.Handler
}

// LimiterConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Limiter], to review the [limiter.Config] that will be configured by [Limiter]
//
// Please see the [limiter] middleware configuration documentation.
//
// [limiter]: https://docs.gofiber.io/middleware/limiter/#config
type LimiterConfig interface {
	LimiterConfig(config *limiter.Config)
}

type limiterConfigurationProperties struct {
	storage                *string
	handler                string
	maxConn                int
	expiration             time.Duration
	disableheaders         bool
	disableValueRedaction  bool
	skipFailedRequests     bool
	skipSuccessfulRequests bool
}

type limiterConfig struct {
	limiter limiterConfigurationProperties
}

var _, _ = any(&FiberMwConfig[limiterConfig]{}).(Stereotype)

func (this *limiterConfig) Config() limiter.Config {
	config := limiter.ConfigDefault

	switch String(this.limiter.handler).ToLowerCase() {
	case "sliding", "slidingwindow":
		config.LimiterMiddleware = limiter.SlidingWindow{}
	case "fixed", "fixedwindow":
		config.LimiterMiddleware = limiter.FixedWindow{}
	default:
		config.LimiterMiddleware = limiter.ConfigDefault.LimiterMiddleware
	}

	config.Max = this.limiter.maxConn
	config.Expiration = this.limiter.expiration
	config.DisableHeaders = this.limiter.disableheaders
	config.DisableValueRedaction = this.limiter.disableValueRedaction
	config.SkipFailedRequests = this.limiter.skipFailedRequests
	config.SkipSuccessfulRequests = this.limiter.skipSuccessfulRequests

	return config
}

// [Limiter] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  Limiter
//	}
//
//	func (this *ApiController) LimiterConfig(config *limiter.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [LimiterNext], [LimiterMaxFunc], [LimiterKeyGenerator] and [LimiterLimitReached].
//
// Also, the parent container may implements [limiterConfig] to receive a pointer to the [limiter.Config]
// that will be used.
//
// Please see those interface definitions and the [limiter] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					limiter:
//						handler: fixed
//						maxConn: 5 # Renamed from max
//						expiration: 1m
//						skipFailedRequests: false
//						skipSuccessfulRequests: false
//
// [limiter]: https://docs.gofiber.io/middleware/limiter/#config
type Limiter struct {
	config          *FiberMwConfig[limiterConfig] `@:"Inject"`
	cfg             *limiter.Config
	storageProvider StorageProvider `@:"Inject"`
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *Limiter) PostConstruct(parent any) Error {
	cfg := this.config.GetConfig(parent)
	this.cfg = new(cfg.Config())
	if cfg.limiter.storage != nil {
		this.cfg.Storage = this.storageProvider.Provide(*cfg.limiter.storage)
	}
	if next, ok := parent.(LimiterNext); ok {
		this.cfg.Next = next.LimiterNext()
	}
	if maxFunc, ok := parent.(LimiterMaxFunc); ok {
		this.cfg.MaxFunc = maxFunc.LimiterMaxFunc()
	}
	if keyGenerator, ok := parent.(LimiterKeyGenerator); ok {
		this.cfg.KeyGenerator = keyGenerator.LimiterKeyGenerator()
	}
	if limitReached, ok := parent.(LimiterLimitReached); ok {
		this.cfg.LimitReached = limitReached.LimiterLimitReached()
	}
	if config, ok := parent.(LimiterConfig); ok {
		config.LimiterConfig(this.cfg)
	}
	return nil
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *Limiter) Handler() fiber.Handler {
	return limiter.New(*this.cfg)
}
