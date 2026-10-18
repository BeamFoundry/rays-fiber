package middleware

import (
	"time"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	. "github.com/BeamFoundry/rays-web-security/pkg/core"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
)

// CsrfNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Csrf], to indicates when to skip it this middleware.
//
// Please see the [csrf] middleware configuration documentation.
//
// [csrf]: https://docs.gofiber.io/middleware/csrf/#config-properties
type CsrfNext interface {
	CsrfNext() func(fiber.Ctx) bool
}

// CsrfExtractor is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Csrf], to indicates a token extraction method.
//
// Please see the [csrf] middleware configuration documentation.
//
// [csrf]: https://docs.gofiber.io/middleware/csrf/#config-properties
type CsrfExtractor interface {
	CsrfExtractor() extractors.Extractor
}

// CsrfKeyGenerator is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Csrf], to indicates a token generation function.
//
// Please see the [csrf] middleware configuration documentation.
//
// [csrf]: https://docs.gofiber.io/middleware/csrf/#config-properties
type CsrfKeyGenerator interface {
	CsrfKeyGenerator() func() string
}

// CsrfErrorHandler is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Csrf], to indicates when to skip it this middleware.
//
// Please see the [csrf] middleware configuration documentation.
//
// [csrf]: https://docs.gofiber.io/middleware/csrf/#config-properties
type CsrfErrorHandler interface {
	CsrfErrorHandler() func(fiber.Ctx, error) error
}

// CsrfConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Csrf], to review the [csrf.Config] that will be configured by [Csrf]
//
// Please see the [csrf] middleware configuration documentation.
//
// [csrf]: https://docs.gofiber.io/middleware/csrf/#config-properties
type CsrfConfig interface {
	CsrfConfig(config *csrf.Config)
}

// CsrfSecurityConfig is an interface to be defined by a component that will customize
// the [Csrf] middleware when an HttpSecurity is defined for that server instance.
//
// Please see the [csrf] middleware configuration documentation.
//
// [csrf]: https://docs.gofiber.io/middleware/csrf/#config-properties
type CsrfSecurityConfig interface {
	CsrfSecurityConfig(server string, config *csrf.Config)
}

type csrfConfigurationProperties struct {
	storage               *string
	cookieName            string
	cookieDomain          string
	cookiePath            string
	cookieSameSite        string
	trustedOrigins        []string
	idleTimeout           time.Duration
	cookieSecure          bool
	cookieHTTPOnly        bool
	cookieSessionOnly     bool
	disableValueRedaction bool
	singleUseToken        bool
}

type csrfConfig struct {
	csrf csrfConfigurationProperties
}

var _, _ = any(&FiberMwConfig[csrfConfig]{}).(Stereotype)

func (this *csrfConfig) Config() csrf.Config {
	config := csrf.ConfigDefault

	config.CookieName = this.csrf.cookieName
	config.CookieDomain = this.csrf.cookieDomain
	config.CookiePath = this.csrf.cookiePath
	config.CookieSameSite = this.csrf.cookieSameSite
	config.TrustedOrigins = this.csrf.trustedOrigins
	config.IdleTimeout = this.csrf.idleTimeout
	config.CookieSecure = this.csrf.cookieSecure
	config.CookieHTTPOnly = this.csrf.cookieHTTPOnly
	config.CookieSessionOnly = this.csrf.cookieSessionOnly
	config.DisableValueRedaction = this.csrf.disableValueRedaction
	config.SingleUseToken = this.csrf.singleUseToken

	return config
}

// [Csrf] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// NOTE: This middleware is marked with the [github.com/BeamFoundry/rays/pkg/core.PostAuthorisationFilter] interface.
// As such it will be skipped when used with [github.com/BeamFoundry/rays-web-security] but it will be configured and
// used automatically, if needed.
//
// WARNING: This middleware may depends on the [Session] middleware. As such the [Session] middleware MUST
// always be configured before [Csrf]. Please see the example below, the [session] and [csrf] Fiber documentation.
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  Session // Session should always be initialized before Csrf.
//	  Csrf
//	}
//
//	func (this *ApiController) CsrfConfig(config *csrf.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [CsrfNext], [CsrfExtractor], [CsrfKeyGenerator] and [CsrfErrorHandler].
//
// Also, the parent container may implements [CsrfConfig] to receive a pointer to the [csrf.Config]
// that will be used.
//
// Please see those interface definitions and the [csrf] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					csrf:
//						cookieName: "csrf_"
//						cookieDomain: ""
//						cookiePath: ""
//						cookieSameSite: "Lax"
//						trustedOrigins: []
//						idleTimeout: 30m
//						cookieSecure: false
//						cookieHTTPOnly: false
//						cookieSessionOnly: false
//						disableValueRedaction: false
//						singleUseToken: false
//
// [csrf]: https://docs.gofiber.io/middleware/csrf/#config-properties
// [session]: https://docs.gofiber.io/middleware/session/#config-properties
type Csrf struct {
	PostAuthorisationFilter
	config *FiberMwConfig[csrfConfig] `@:"Inject"`
	cfg    *csrf.Config

	storageProvider StorageProvider `@:"Inject"`
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *Csrf) PostConstruct(parent any) Error {
	cfg := this.config.GetConfig(parent)
	this.cfg = new(cfg.Config())
	if cfg.csrf.storage != nil {
		this.cfg.Storage = this.storageProvider.Provide(*cfg.csrf.storage)
	}
	if next, ok := parent.(CsrfNext); ok {
		this.cfg.Next = next.CsrfNext()
	}
	// SessionStore is implemented by [Session]
	if session, ok := parent.(SessionStore); ok {
		store := session.SessionStore()
		if store == nil {
			return NewError("IllegalArgumentException", "CSRF SessionStore is nil. Have you initialized Session before Csrf ?")
		}
		this.cfg.Session = store
	}
	if extractor, ok := parent.(CsrfExtractor); ok {
		this.cfg.Extractor = extractor.CsrfExtractor()
	}
	if errorHandler, ok := parent.(CsrfErrorHandler); ok {
		this.cfg.ErrorHandler = errorHandler.CsrfErrorHandler()
	}
	if keyGenerator, ok := parent.(CsrfKeyGenerator); ok {
		this.cfg.KeyGenerator = keyGenerator.CsrfKeyGenerator()
	}
	if config, ok := parent.(CsrfConfig); ok {
		config.CsrfConfig(this.cfg)
	}
	return nil
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *Csrf) Handler() fiber.Handler {
	return csrf.New(*this.cfg)
}
