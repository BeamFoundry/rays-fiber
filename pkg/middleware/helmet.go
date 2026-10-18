package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/helmet"
)

// HelmetNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Helmet], to indicates when to skip it this middleware.
//
// Please see the [helmet] middleware configuration documentation.
//
// [helmet]: https://docs.gofiber.io/middleware/helmet/#config
type HelmetNext interface {
	HelmetNext() func(fiber.Ctx) bool
}

// HelmetConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Helmet], to review the [helmet.Config] that will be configured by [Helmet]
//
// Please see the [helmet] middleware configuration documentation.
//
// [helmet]: https://docs.gofiber.io/middleware/helmet/#config
type HelmetConfig interface {
	HelmetConfig(config *helmet.Config)
}

type helmetConfigurationProperties struct {
	xSSProtection             string
	contentTypeNosniff        string
	xFrameOptions             string
	contentSecurityPolicy     string
	referrerPolicy            string
	permissionPolicy          string
	crossOriginEmbedderPolicy string
	crossOriginOpenerPolicy   string
	crossOriginResourcePolicy string
	originAgentCluster        string
	xDNSPrefetchControl       string
	xDownloadOptions          string
	xPermittedCrossDomain     string
	hstsMaxAge                int
	hstsExcludeSubdomains     bool
	cspReportOnly             bool
	hstsPreloadEnabled        bool
}

type helmetConfig struct {
	helmet helmetConfigurationProperties
}

var _, _ = any(&FiberMwConfig[helmetConfig]{}).(Stereotype)

func (this *helmetConfig) Config() helmet.Config {
	config := helmet.ConfigDefault

	config.XSSProtection = this.helmet.xSSProtection
	config.ContentTypeNosniff = this.helmet.contentTypeNosniff
	config.XFrameOptions = this.helmet.xFrameOptions
	config.ContentSecurityPolicy = this.helmet.contentSecurityPolicy
	config.ReferrerPolicy = this.helmet.referrerPolicy
	config.PermissionPolicy = this.helmet.permissionPolicy
	config.CrossOriginEmbedderPolicy = this.helmet.crossOriginEmbedderPolicy
	config.CrossOriginOpenerPolicy = this.helmet.crossOriginOpenerPolicy
	config.CrossOriginResourcePolicy = this.helmet.crossOriginResourcePolicy
	config.OriginAgentCluster = this.helmet.originAgentCluster
	config.XDNSPrefetchControl = this.helmet.xDNSPrefetchControl
	config.XDownloadOptions = this.helmet.xDownloadOptions
	config.XPermittedCrossDomain = this.helmet.xPermittedCrossDomain
	config.HSTSMaxAge = this.helmet.hstsMaxAge
	config.HSTSExcludeSubdomains = this.helmet.hstsExcludeSubdomains
	config.CSPReportOnly = this.helmet.cspReportOnly
	config.HSTSPreloadEnabled = this.helmet.hstsPreloadEnabled

	return config
}

// [Helmet] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  Helmet
//	}
//
//	func (this *ApiController) HelmetConfig(config *helmet.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through the interface [HelmetNext] the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
//
// Also, the parent container may implements [HelmetConfig] to receive a pointer to the [helmet.Config]
// that will be used.
//
// Please see those interface definitions and the [helmet] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					helmet:
//						xSSProtection: "0"
//						contentTypeNosniff: "nosniff"
//						xFrameOptions: "SAMEORIGIN"
//						contentSecurityPolicy: ""
//						referrerPolicy: "no-referrer"
//						permissionPolicy: ""
//						crossOriginEmbedderPolicy: "require-corp"
//						crossOriginOpenerPolicy: "same-origin"
//						crossOriginResourcePolicy: "same-origin"
//						originAgentCluster: "?1"
//						xDNSPrefetchControl: "off"
//						xDownloadOptions: "noopen"
//						xPermittedCrossDomain: "none"
//						hstsMaxAge: 0
//						hstsExcludeSubdomains: false
//						cspReportOnly: false
//						hstsPreloadEnabled: false
//
// [helmet]: https://docs.gofiber.io/middleware/helmet/#config
type Helmet struct {
	config *FiberMwConfig[helmetConfig] `@:"Inject"`
	cfg    *helmet.Config
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *Helmet) PostConstruct(parent any) Error {
	this.cfg = new(this.config.GetConfig(parent).Config())
	if next, ok := parent.(HelmetNext); ok {
		this.cfg.Next = next.HelmetNext()
	}
	if config, ok := parent.(HelmetConfig); ok {
		config.HelmetConfig(this.cfg)
	}
	return nil
}

func (this *Helmet) Handler() fiber.Handler {
	return helmet.New(*this.cfg)
}
