package middleware

import (
	"io/fs"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/favicon"
)

// FaviconNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Favicon], to indicates when to skip it this middleware.
//
// Please see the [favicon] middleware configuration documentation.
//
// [favicon]: https://docs.gofiber.io/middleware/favicon/#config
type FaviconNext interface {
	FaviconNext() func(fiber.Ctx) bool
}

// FaviconFileSystem is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Favicon], to indicates an optional alternate filesystem from which to load the favicon file.
//
// Please see the [favicon] middleware configuration documentation.
//
// [favicon]: https://docs.gofiber.io/middleware/favicon/#config
type FaviconFileSystem interface {
	FaviconFileSystem() fs.FS
}

// FaviconConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Favicon], to review the [favicon.Config] that will be configured by [Favicon]
//
// Please see the [favicon] middleware configuration documentation.
//
// [favicon]: https://docs.gofiber.io/middleware/favicon/#config
type FaviconConfig interface {
	FaviconConfig(config *favicon.Config)
}

type faviconConfigurationProperties struct {
	file         string
	url          string
	cacheControl string
	maxBytes     int64
	// TODO: Handle Data ([]byte) or not ???
}

type faviconConfig struct {
	favicon faviconConfigurationProperties
}

var _, _ = any(&FiberMwConfig[faviconConfig]{}).(Stereotype)

func (this *faviconConfig) Config() favicon.Config {
	config := favicon.ConfigDefault

	config.File = this.favicon.file
	config.URL = this.favicon.url
	config.CacheControl = this.favicon.cacheControl
	config.MaxBytes = this.favicon.maxBytes

	return config
}

// [Favicon] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  Favicon
//	}
//
//	func (this *ApiController) FaviconConfig(config *favicon.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [FaviconNext] and [FaviconFileSystem].
//
// Also, the parent container may implements [FaviconConfig] to receive a pointer to the [favicon.Config]
// that will be used.
//
// Please see those interface definitions and the [favicon] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					favicon:
//						file: ""
//						url: "/favicon.ico"
//						cacheControl: "public, max-age=31536000"
//
// [favicon]: https://docs.gofiber.io/middleware/favicon/#config
type Favicon struct {
	config *FiberMwConfig[faviconConfig] `@:"Inject"`
	cfg    *favicon.Config
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *Favicon) PostConstruct(parent any) Error {
	this.cfg = new(this.config.GetConfig(parent).Config())
	if next, ok := parent.(FaviconNext); ok {
		this.cfg.Next = next.FaviconNext()
	}
	if fileSystem, ok := parent.(FaviconFileSystem); ok {
		this.cfg.FileSystem = fileSystem.FaviconFileSystem()
	}
	return nil
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *Favicon) Handler() fiber.Handler {
	return favicon.New(*this.cfg)
}
