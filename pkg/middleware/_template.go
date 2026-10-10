package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/BeamFoundry/rays-fiber/pkg/storage"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/template"
)

// TemplateNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Template], to indicates when to skip it this middleware.
//
// Please see the [template] middleware configuration documentation.
//
// [template]: https://docs.gofiber.io/middleware/template/#config
type TemplateNext interface {
	TemplateNext() func(fiber.Ctx) bool
}

// TemplateConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Template], to review the [template.Config] that will be configured by [template]
//
// Please see the [template] middleware configuration documentation.
//
// [template]: https://docs.gofiber.io/middleware/template/#config
type TemplateConfig interface {
	TemplateConfig(config *template.Config)
}

type templateConfigurationProperties struct {
	realm string
}

type templateConfig struct {
	template templateConfigurationProperties
}

var _, _ = any(&FiberMwConfig[templateConfig]{}).(Stereotype)

func (this *templateConfig) Config() template.Config {
	config := template.ConfigDefault

	config.Realm = this.template.realm

	return config
}

// [Template] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//		FiberController
//		RequestMapping `Path:"/api"`
//		Template
//	}
//
//	func (this *ApiController) TemplateConfig(config *template.Config) {
//		// changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [TemplateNext], [], [] and [].
//
// Also, the parent container may implements [TemplateConfig] to receive a pointer to the [template.Config]
// that will be used.
//
// Please see those interface definitions and the [template] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//	        template:
//	          expiration: 5m
//
// [template]: https://docs.gofiber.io/middleware/template/#config
type Template struct {
	config          *FiberMwConfig[templateConfig] `@:"Inject"`
	storageProvider StorageProvider                `@:"Inject"`
	cfg             *template.Config
}

func (this *Template) PostConstruct(parent any) Error {
	cfg := this.config.GetConfig(parent)
	this.cfg = new(cfg.Config())
	if cfg.template.storage != nil {
		this.cfg.Storage = this.storageProvider.Provide(*cfg.template.storage)
	}
	if next, ok := (parent).(TemplateNext); ok {
		this.cfg.Next = next.TemplateNext()
	}
	if config, ok := (parent).(TemplateConfig); ok {
		config.TemplateConfig(this.cfg)
	}
	return nil
}

func (this *Template) Handler() fiber.Handler {
	return template.New(*this.cfg)
}
