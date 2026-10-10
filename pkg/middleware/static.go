package middleware

import (
	"io/fs"
	"os"
	"reflect"
	"time"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

// StaticNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Static], to indicates when to skip it this middleware.
//
// Please see the [static] middleware configuration documentation.
//
// [static]: https://docs.gofiber.io/middleware/static/#config
type StaticNext interface {
	StaticNext() func(fiber.Ctx) bool
}

// StaticFS is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Static], to indicates the file system to serve the static files from.
//
// Please see the [static] middleware configuration documentation.
//
// [static]: https://docs.gofiber.io/middleware/static/#config
type StaticFS interface {
	StaticFS() fs.FS
}

// StaticModifyResponse is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Static], to indicates a function that allows you to alter the response.
//
// Please see the [static] middleware configuration documentation.
//
// [static]: https://docs.gofiber.io/middleware/static/#config
type StaticModifyResponse interface {
	StaticModifyResponse() fiber.Handler
}

// StaticNotFoundHandler is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Static], to indicates a function to handle when the path is not found.
//
// Please see the [static] middleware configuration documentation.
//
// [static]: https://docs.gofiber.io/middleware/static/#config
type StaticNotFoundHandler interface {
	StaticNotFoundHandler() fiber.Handler
}

// StaticConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Static], to review the [static.Config] that will be configured by [Static]
//
// Please see the [static] middleware configuration documentation.
//
// [static]: https://docs.gofiber.io/middleware/static/#config
type StaticConfig interface {
	StaticConfig(config *static.Config)
}

type staticConfigurationProperties struct {
	indexNames    []string
	cacheDuration time.Duration
	maxAge        int
	compress      bool
	byteRange     bool
	browse        bool
	download      bool
}

type staticConfig struct {
	static staticConfigurationProperties
}

var _, _ = any(&FiberMwConfig[staticConfig]{}).(Stereotype)

func (this *staticConfig) Config() static.Config {
	config := static.ConfigDefault

	config.IndexNames = this.static.indexNames
	config.CacheDuration = this.static.cacheDuration
	config.MaxAge = this.static.maxAge
	config.Compress = this.static.compress
	config.ByteRange = this.static.byteRange
	config.Browse = this.static.browse
	config.Download = this.static.download

	return config
}

// [Static] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  Static
//	}
//
//	func (this *ApiController) StaticConfig(config *static.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [StaticNext], [StaticFS], [StaticModifyResponse] and [StaticNotFoundHandler].
//
// Also, the parent container may implements [StaticConfig] to receive a pointer to the [static.Config]
// that will be used.
//
// Please see those interface definitions and the [static] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					static:
//						indexNames: ["index.html"]
//						cacheDuration: 10s
//						maxAge: 0
//						compress: false
//						byteRange: false
//						browse: false
//						download: false
//
// [static]: https://docs.gofiber.io/middleware/static/#config
type Static struct {
	config *FiberMwConfig[staticConfig] `@:"Inject"`
	cfg    *static.Config

	staticFS fs.FS
}

func (this *Static) PostConstruct(parent any) Error {
	reflectType := reflect.TypeOf(parent).Elem()
	this.staticFS = nil

	// TODO: Remove tag
	if staticFS, ok := (parent).(StaticFS); ok {
		this.staticFS = staticFS.StaticFS()
	} else if field, ok := reflectType.FieldByName("Static"); ok {
		if dir, ok := field.Tag.Lookup("Dir"); ok {
			this.staticFS = os.DirFS(dir)
		}
	}

	this.cfg = new(this.config.GetConfig(parent).Config())
	if this.staticFS != nil {
		this.cfg.FS = this.staticFS
	}
	if next, ok := (parent).(StaticNext); ok {
		this.cfg.Next = next.StaticNext()
	}
	if modifyResponse, ok := (parent).(StaticModifyResponse); ok {
		this.cfg.ModifyResponse = modifyResponse.StaticModifyResponse()
	}
	if notFoundHandler, ok := (parent).(StaticNotFoundHandler); ok {
		this.cfg.NotFoundHandler = notFoundHandler.StaticNotFoundHandler()
	}
	if config, ok := (parent).(StaticConfig); ok {
		config.StaticConfig(this.cfg)
	}
	return nil
}

func (this *Static) Handler() fiber.Handler {
	return static.New("", *this.cfg)
}
