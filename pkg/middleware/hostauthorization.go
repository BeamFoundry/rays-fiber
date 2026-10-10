package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/hostauthorization"
)

// HostAuthorizationNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [HostAuthorization], to indicates when to skip it this middleware.
//
// Please see the [hostauthorization] middleware configuration documentation.
//
// [hostauthorization]: https://docs.gofiber.io/middleware/hostauthorization/#config
type HostAuthorizationNext interface {
	HostAuthorizationNext() func(fiber.Ctx) bool
}

// HostAuthorizationAllowedHostsFunc is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [HostAuthorization], to indicates dynamic validator called only when no static AllowedHosts rule matches.
//
// Please see the [hostauthorization] middleware configuration documentation.
//
// [hostauthorization]: https://docs.gofiber.io/middleware/hostauthorization/#config
type HostAuthorizationAllowedHostsFunc interface {
	HostAuthorizationAllowedHostsFunc() func(host string) bool
}

// HostAuthorizationNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [HostAuthorization], to indicates an error handler called when a request is rejected.
//
// Please see the [hostauthorization] middleware configuration documentation.
//
// [hostauthorization]: https://docs.gofiber.io/middleware/hostauthorization/#config
type HostAuthorizationErrorHandler interface {
	HostAuthorizationErrorHandler() fiber.ErrorHandler
}

// HostAuthorizationConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [HostAuthorization], to review the [hostauthorization.Config] that will be configured by [hostauthorization]
//
// Please see the [hostauthorization] middleware configuration documentation.
//
// [hostauthorization]: https://docs.gofiber.io/middleware/hostauthorization/#config
type HostAuthorizationConfig interface {
	HostAuthorizationConfig(config *hostauthorization.Config)
}

type hostauthorizationConfigurationProperties struct {
	allowedHosts []string
}

type hostauthorizationConfig struct {
	hostauthorization hostauthorizationConfigurationProperties
}

var _, _ = any(&FiberMwConfig[hostauthorizationConfig]{}).(Stereotype)

func (this *hostauthorizationConfig) Config() hostauthorization.Config {
	config := hostauthorization.ConfigDefault

	config.AllowedHosts = this.hostauthorization.allowedHosts

	return config
}

// [HostAuthorization] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//		FiberController
//		RequestMapping `Path:"/api"`
//		HostAuthorization
//	}
//
//	func (this *ApiController) HostAuthorizationConfig(config *hostauthorization.Config) {
//		// changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [HostAuthorizationNext], [], [] and [].
//
// Also, the parent container may implements [HostAuthorizationConfig] to receive a pointer to the [hostauthorization.Config]
// that will be used.
//
// Please see those interface definitions and the [hostauthorization] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//	        hostauthorization:
//	          allowedHosts: []
//
// [hostauthorization]: https://docs.gofiber.io/middleware/hostauthorization/#config
type HostAuthorization struct {
	config *FiberMwConfig[hostauthorizationConfig] `@:"Inject"`
	cfg    *hostauthorization.Config
}

func (this *HostAuthorization) PostConstruct(parent any) Error {
	this.cfg = new(this.config.GetConfig(parent).Config())
	if next, ok := (parent).(HostAuthorizationNext); ok {
		this.cfg.Next = next.HostAuthorizationNext()
	}
	if allowedHosts, ok := (parent).(HostAuthorizationAllowedHostsFunc); ok {
		this.cfg.AllowedHostsFunc = allowedHosts.HostAuthorizationAllowedHostsFunc()
	}
	if errorHandler, ok := (parent).(HostAuthorizationErrorHandler); ok {
		this.cfg.ErrorHandler = errorHandler.HostAuthorizationErrorHandler()
	}
	if config, ok := (parent).(HostAuthorizationConfig); ok {
		config.HostAuthorizationConfig(this.cfg)
	}
	return nil
}

func (this *HostAuthorization) Handler() fiber.Handler {
	return hostauthorization.New(*this.cfg)
}
