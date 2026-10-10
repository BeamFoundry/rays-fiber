package middleware

import (
	"reflect"
	"time"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/timeout"
)

// TimeoutNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Timeout], to indicates when to skip it this middleware.
//
// Please see the [timeout] middleware configuration documentation.
//
// [timeout]: https://docs.gofiber.io/middleware/timeout/#config
type TimeoutNext interface {
	TimeoutNext() func(fiber.Ctx) bool
}

// TimeoutOnTimeout is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Timeout], to indicates a handler executed when a timeout occurs.
//
// Please see the [timeout] middleware configuration documentation.
//
// [timeout]: https://docs.gofiber.io/middleware/timeout/#config
type TimeoutOnTimeout interface {
	TimeoutOnTimeout() fiber.Handler
}

// TimeoutErrors is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Timeout], to indicates custom errors treated as timeout errors.
//
// Please see the [timeout] middleware configuration documentation.
//
// [timeout]: https://docs.gofiber.io/middleware/timeout/#config
type TimeoutErrors interface {
	TimeoutErrors() []error
}

// TimeoutConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Timeout], to review the [timeout.Config] that will be configured by [Timeout]
//
// Please see the [timeout] middleware configuration documentation.
//
// [timeout]: https://docs.gofiber.io/middleware/timeout/#config
type TimeoutConfig interface {
	TimeoutConfig(config *timeout.Config)
}

type timeoutConfigurationProperties struct {
	timeout time.Duration
}

func (this *timeoutConfigurationProperties) Timeout() time.Duration {
	return this.timeout
}

type timeoutConfig struct {
	timeout timeoutConfigurationProperties
}

var _, _ = any(&FiberMwConfig[timeoutConfig]{}).(Stereotype)

func (this *timeoutConfig) Config() timeout.Config {
	config := timeout.Config{}

	config.Timeout = this.timeout.Timeout()

	return config
}

// [Timeout] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  Timeout `Duration:"10s"`
//	}
//
//	func (this *ApiController) TimeoutConfig(config *timeout.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [TimeoutNext], [TimeoutOnTimeout] and [TimeoutErrors].
//
// Also, the parent container may implements [TimeoutConfig] to receive a pointer to the [timeout.Config]
// that will be used.
//
// You also have the possibility to define the timeout duration as a struct tag. See the example above.
//
// Please see those interface definitions and the [timeout] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					timeout:
//						timeout: 10s
//
// [timeout]: https://docs.gofiber.io/middleware/timeout/#config
//
//nolint:dupword // Needed
type Timeout struct {
	config *FiberMwConfig[timeoutConfig] `@:"Inject"`
	cfg    *timeout.Config
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *Timeout) PostConstruct(parent any) Error {
	this.cfg = new(this.config.GetConfig(parent).Config())
	// TODO: Remove tag
	reflectType := reflect.TypeOf(parent).Elem()
	if field, ok := reflectType.FieldByName("Timeout"); ok {
		if duration, ok := field.Tag.Lookup("Duration"); ok {
			if d, err := time.ParseDuration(duration); err == nil {
				this.cfg.Timeout = d
			}
		}
	}
	if next, ok := (parent).(TimeoutNext); ok {
		this.cfg.Next = next.TimeoutNext()
	}

	if errors, ok := (parent).(TimeoutErrors); ok {
		this.cfg.Errors = errors.TimeoutErrors()
	}

	if config, ok := (parent).(TimeoutConfig); ok {
		config.TimeoutConfig(this.cfg)
	}
	return nil
}

// Filter() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberFilter]
func (this *Timeout) Filter(handler fiber.Handler) fiber.Handler {
	return timeout.New(handler, *this.cfg)
}
