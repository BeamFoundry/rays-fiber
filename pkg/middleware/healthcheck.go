package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
)

// HealthCheckNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [HealthCheck], to indicates when to skip it this middleware.
//
// Please see the [healthcheck] middleware configuration documentation.
//
// [healthcheck]: https://docs.gofiber.io/middleware/healthcheck/#config
type HealthCheckNext interface {
	HealthCheckNext() func(fiber.Ctx) bool
}

// HealthCheckProbe is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [HealthCheck], to indicates a function that will be used to check probe the health.
//
// Please see the [healthcheck] middleware configuration documentation.
//
// [healthcheck]: https://docs.gofiber.io/middleware/healthcheck/#config
type HealthCheckProbe interface {
	HealthCheckProbe() func(fiber.Ctx) bool
}

// HealthCheckConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [HealthCheck], to review the [healthcheck.Config] that will be configured by [HealthCheck]
//
// Please see the [healthcheck] middleware configuration documentation.
//
// [healthcheck]: https://docs.gofiber.io/middleware/healthcheck/#config
type HealthCheckConfig interface {
	HealthCheckConfig(config *healthcheck.Config)
}

type healthCheckConfigurationProperties struct {
	format HealthCheckResponseFormat
}

type healthCheckConfig struct {
	healthcheck healthCheckConfigurationProperties
}

// HealthCheckResponseFormat is an Enum type to convert format to use from a string.
// Please see the [healthcheck] middleware configuration documentation.
//
// [healthcheck]: https://docs.gofiber.io/middleware/healthcheck/#config
type HealthCheckResponseFormat struct {
	Enum[HealthCheckResponseFormat]
	TEXT,
	JSON,
	XML,
	MSGPACK,
	CBOR Enum[HealthCheckResponseFormat]
}

var _, _ = any(&FiberMwConfig[healthCheckConfig]{}).(Stereotype)

func (this *healthCheckConfig) Config() healthcheck.Config {
	config := healthcheck.Config{}
	config.ResponseFormat = healthcheck.ResponseFormat(this.healthcheck.format.Ordinal())
	return config
}

// [HealthCheck] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  HealthCheck
//	}
//
//	func (this *ApiController) HealthCheckConfig(config *healthcheck.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [HealthCheckNext] and [HealthCheckProbe].
//
// Also, the parent container may implements [HealthCheckConfig] to receive a pointer to the [healthcheck.Config]
// that will be used.
//
// Please see those interface definitions and the [healthcheck] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					healthcheck:
//						format: "TEXT" # Possible values : TEXT, JSON, XML, MSGPACK, CBOR
//
// [healthcheck]: https://docs.gofiber.io/middleware/healthcheck/#config
type HealthCheck struct {
	config *FiberMwConfig[healthCheckConfig] `@:"Inject"`
	cfg    *healthcheck.Config
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *HealthCheck) PostConstruct(parent any) Error {
	this.cfg = new(this.config.GetConfig(parent).Config())
	if next, ok := (parent).(HealthCheckNext); ok {
		this.cfg.Next = next.HealthCheckNext()
	}
	if probe, ok := (parent).(HealthCheckProbe); ok {
		this.cfg.Probe = probe.HealthCheckProbe()
	}
	return nil
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *HealthCheck) Handler() fiber.Handler {
	return healthcheck.New(*this.cfg)
}
