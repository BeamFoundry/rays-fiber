package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
)

// CompressNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Compress], to indicates when to skip it this middleware.
//
// Please see the [compress] middleware configuration documentation.
//
// [compress]: https://docs.gofiber.io/middleware/compress/#config
type CompressNext interface {
	CompressNext() func(fiber.Ctx) bool
}

// CompressConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Compress], to review the [compress.Config] that will be configured by [Compress]
//
// Please see the [compress] middleware configuration documentation.
//
// [compress]: https://docs.gofiber.io/middleware/compress/#config
type CompressConfig interface {
	CompressConfig(config *compress.Config)
}

type compressConfigurationProperties struct {
	level string
}

type compressConfig struct {
	compress compressConfigurationProperties
}

var _, _ = any(&FiberMwConfig[compressConfig]{}).(Stereotype)

func (this *compressConfig) Config() compress.Config {
	config := compress.ConfigDefault
	switch String(this.compress.level).ToUpperCase().String() {
	case "DISABLED":
		config.Level = compress.LevelDisabled
	case "BESTSPEED", "SPEED":
		config.Level = compress.LevelBestSpeed
	case "BESTCOMPRESSION", "COMPRESSION":
		config.Level = compress.LevelBestCompression
	default:
		config.Level = compress.LevelDefault
	}
	return config
}

// [Compress] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//		FiberController
//		RequestMapping `Path:"/api"`
//		Compress
//	}
//
//	func (this *ApiController) CompressConfig(config *compress.Config) {
//		// changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [CompressNext]. Alternatively, the parent container may implements [CompressConfig] to receive a pointer
// to the [compress.Config] that will be used.
//
// Please see those interface definitions and the [compress] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					compress:
//						level: "default" # Case Unsensitive. Possible values : Default, Disabled, Speed, BestSpeed, Compression, BestCompression.
//
// [compress]: https://docs.gofiber.io/middleware/compress/#config
type Compress struct {
	config *FiberMwConfig[compressConfig] `@:"Inject"`
	cfg    *compress.Config
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *Compress) PostConstruct(parent any) Error {
	this.cfg = new(this.config.GetConfig(parent).Config())
	if next, ok := parent.(CompressNext); ok {
		this.cfg.Next = next.CompressNext()
	}
	if config, ok := parent.(CompressConfig); ok {
		config.CompressConfig(this.cfg)
	}
	return nil
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *Compress) Handler() fiber.Handler {
	return compress.New(*this.cfg)
}
