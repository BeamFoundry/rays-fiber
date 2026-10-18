package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/lang"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/earlydata"
)

// EarlyDataNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [EarlyData], to indicates when to skip it this middleware.
//
// Please see the [earlydata] middleware configuration documentation.
//
// [earlydata]: https://docs.gofiber.io/middleware/earlydata/#config
type EarlyDataNext interface {
	EarlyDataNext() func(fiber.Ctx) bool
}

// EarlyDataIsEarlyData is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [EarlyData], to indicates whether the request used early data.
//
// Please see the [earlydata] middleware configuration documentation.
//
// [earlydata]: https://docs.gofiber.io/middleware/earlydata/#config
type EarlyDataIsEarlyData interface {
	EarlyDataIsEarlyData() func(c fiber.Ctx) bool
}

// EarlyDataAllowEarlyData is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [EarlyData], to indicates if an early-data request should be allowed.
//
// Please see the [earlydata] middleware configuration documentation.
//
// [earlydata]: https://docs.gofiber.io/middleware/earlydata/#config
type EarlyDataAllowEarlyData interface {
	EarlyDataAllowEarlyData() func(c fiber.Ctx) bool
}

// EarlyDataError is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [EarlyData], to indicates which error to return when an early-data request is rejected.
// Defaults to [fiber.ErrTooEarly].
//
// Please see the [earlydata] middleware configuration documentation.
//
// [earlydata]: https://docs.gofiber.io/middleware/earlydata/#config
type EarlyDataError interface {
	EarlyDataError() error
}

// EarlyDataConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [EarlyData], to review the [earlydata.Config] that will be configured by [EarlyData]
//
// Please see the [earlydata] middleware configuration documentation.
//
// [earlydata]: https://docs.gofiber.io/middleware/earlydata/#config
type EarlyDataConfig interface {
	EarlyDataConfig(config *earlydata.Config)
}

// [EarlyData] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  EarlyData
//	}
//
//	func (this *ApiController) EarlyDataConfig(config *earlydata.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [EarlyDataNext], [EarlyDataIsEarlyData], [EarlyDataAllowEarlyData] and [EarlyDataError].
//
// Also, the parent container may implements [EarlyDataConfig] to receive a pointer to the [earlydata.Config]
// that will be used.
//
// Please see those interface definitions and the [earlydata] middleware documentation.
//
// [earlydata]: https://docs.gofiber.io/middleware/earlydata/#config
type EarlyData struct {
	cfg *earlydata.Config
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *EarlyData) PostConstruct(parent any) Error {
	this.cfg = new(earlydata.ConfigDefault)
	if next, ok := parent.(EarlyDataNext); ok {
		this.cfg.Next = next.EarlyDataNext()
	}
	if isEarlyData, ok := parent.(EarlyDataIsEarlyData); ok {
		this.cfg.IsEarlyData = isEarlyData.EarlyDataIsEarlyData()
	}
	if allowEarlyData, ok := parent.(EarlyDataAllowEarlyData); ok {
		this.cfg.AllowEarlyData = allowEarlyData.EarlyDataAllowEarlyData()
	}
	if earlyDataError, ok := parent.(EarlyDataError); ok {
		this.cfg.Error = earlyDataError.EarlyDataError()
	}
	if config, ok := parent.(EarlyDataConfig); ok {
		config.EarlyDataConfig(this.cfg)
	}
	return nil
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *EarlyData) Handler() fiber.Handler {
	return earlydata.New(*this.cfg)
}
