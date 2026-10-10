package middleware

import (
	"time"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/idempotency"
)

// IdempotencyNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Idempotency], to indicates when to skip it this middleware.
//
// Please see the [idempotency] middleware configuration documentation.
//
// [idempotency]: https://docs.gofiber.io/middleware/idempotency/#config
type IdempotencyNext interface {
	IdempotencyNext() func(fiber.Ctx) bool
}

// IdempotencyKeyHeaderValidate is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Idempotency], to indicates a function to validate idempotency header syntax.
//
// Please see the [idempotency] middleware configuration documentation.
//
// [idempotency]: https://docs.gofiber.io/middleware/idempotency/#config
type IdempotencyKeyHeaderValidate interface {
	IdempotencyKeyHeaderValidate() func(string) error
}

// IdempotencyLocker is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Idempotency], to indicates a lock to prevent race conditions..
//
// Please see the [idempotency] middleware configuration documentation.
//
// [idempotency]: https://docs.gofiber.io/middleware/idempotency/#config
type IdempotencyLocker interface {
	IdempotencyLocker() idempotency.Locker
}

// IdempotencyConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Idempotency], to review the [idempotency.Config] that will be configured by [Idempotency]
//
// Please see the [idempotency] middleware configuration documentation.
//
// [idempotency]: https://docs.gofiber.io/middleware/idempotency/#config
type IdempotencyConfig interface {
	IdempotencyConfig(config *idempotency.Config)
}

type idempotencyConfigurationProperties struct {
	storage               *string
	keyHeader             string
	keepResponseHeaders   []string
	lifetime              time.Duration
	disableValueRedaction bool
}

type idempotencyConfig struct {
	idempotency idempotencyConfigurationProperties
}

var _, _ = any(&FiberMwConfig[idempotencyConfig]{}).(Stereotype)

func (this *idempotencyConfig) Config() idempotency.Config {
	config := idempotency.ConfigDefault
	config.Lifetime = this.idempotency.lifetime
	config.KeyHeader = this.idempotency.keyHeader
	if len(this.idempotency.keepResponseHeaders) > 0 {
		config.KeepResponseHeaders = this.idempotency.keepResponseHeaders
	}
	config.DisableValueRedaction = this.idempotency.disableValueRedaction
	return config
}

// [Idempotency] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  Idempotency
//	}
//
//	func (this *ApiController) IdempotencyConfig(config *idempotency.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [IdempotencyNext], [IdempotencyKeyHeaderValidate] and [IdempotencyLocker].
//
// Also, the parent container may implements [IdempotencyConfig] to receive a pointer to the [idempotency.Config]
// that will be used.
//
// Please see those interface definitions and the [idempotency] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					idempotency:
//						lifetime: 30m
//						keyHeader: "X-Idempotency-Key"
//						keepResponseHeaders: []
//						storage: null
//
// [idempotency]: https://docs.gofiber.io/middleware/idempotency/#config
type Idempotency struct {
	config          *FiberMwConfig[idempotencyConfig] `@:"Inject"`
	cfg             *idempotency.Config
	storageProvider StorageProvider `@:"Inject"`
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *Idempotency) PostConstruct(parent any) Error {
	cfg := this.config.GetConfig(parent)
	this.cfg = new(cfg.Config())

	if cfg.idempotency.storage != nil {
		this.cfg.Storage = this.storageProvider.Provide(*cfg.idempotency.storage)
	}
	if next, ok := (parent).(IdempotencyNext); ok {
		this.cfg.Next = next.IdempotencyNext()
	}
	if keyHeaderValidate, ok := (parent).(IdempotencyKeyHeaderValidate); ok {
		this.cfg.KeyHeaderValidate = keyHeaderValidate.IdempotencyKeyHeaderValidate()
	}
	if locker, ok := (parent).(IdempotencyLocker); ok {
		this.cfg.Lock = locker.IdempotencyLocker()
	}
	if config, ok := (parent).(IdempotencyConfig); ok {
		config.IdempotencyConfig(this.cfg)
	}
	return nil
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *Idempotency) Handler() fiber.Handler {
	return idempotency.New(*this.cfg)
}
